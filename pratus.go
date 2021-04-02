package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getPRState(owner string, repo string, number int, token string) (state string, statuses []github.RepoStatus, err error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	g := github.NewClient(tc)

	pr, _, err := g.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return "", nil, err
	}

	commit := pr.GetHead().GetSHA()

	status, _, err := g.Repositories.GetCombinedStatus(ctx, owner, repo, commit, nil)
	if err != nil {
		return "", nil, err
	}

	state = status.GetState()
	statuses = status.Statuses
	return state, statuses, nil
}

func stillPending(statuses []github.RepoStatus) (pending bool) {
	for _, status := range statuses {
		if status.GetState() == "pending" {
			return true
		}
	}
	return false
}

func getFailedURLs(statuses []github.RepoStatus) (failed []string) {
	for _, status := range statuses {
		if status.GetState() != "success" {
			failed = append(failed, status.GetTargetURL())
		}
	}
	return failed
}

func parseGitHubURL(baseURL string, URL string) (owner string, repo string, number int) {
	strippedURL := strings.Replace(URL, baseURL, "", 1)
	splitURL := strings.Split(strippedURL, "/")
	owner = splitURL[0]
	repo = splitURL[1]
	number, _ = strconv.Atoi(splitURL[3])
	return
}

func main() {

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN must be set, aborting.")
		os.Exit(1)
	}

	URL := os.Args[1]
	gitHubBaseURL := "https://github.com/"

	timer := time.Duration(60)
	if _, ok := os.LookupEnv("PRATUS_SLEEP_TIMER"); ok {
		customTimer, _ := strconv.Atoi(os.Getenv("PRATUS_SLEEP_TIMER"))
		timer = time.Duration(customTimer)
	}

	sleepTimer := timer * time.Second

	owner, repo, number := parseGitHubURL(gitHubBaseURL, URL)

	fmt.Printf("Checking status of pull request %d in %s/%s every %d seconds\n", number, owner, repo, timer)

	for {

		state, statuses, err := getPRState(owner, repo, number, token)
		if err != nil {
			fmt.Println(err)
			time.Sleep(sleepTimer)
			continue
		}

		if state == "pending" || stillPending(statuses) {
			fmt.Print(".")
			time.Sleep(sleepTimer)
			continue
		}

		failedURLs := getFailedURLs(statuses)

		switch state {
		case "success":
			fmt.Println("\nPR succeeded :)")
			os.Exit(0)
		case "failure", "error":
			fmt.Println("\nPR failed :( Failed URLs:")
			fmt.Println(strings.Join(failedURLs, "\n"))
			os.Exit(1)
		default:
			fmt.Printf("\nPR contained an unknown status: %q", state)
			os.Exit(1)
		}
	}
}
