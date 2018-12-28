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

func getPRState(owner string, repo string, number int) (state string, err error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	g := github.NewClient(tc)

	pr, _, err := g.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return "", err
	}

	commit := pr.GetHead().GetSHA()

	status, _, err := g.Repositories.GetCombinedStatus(ctx, owner, repo, commit, nil)
	if err != nil {
		return "", err
	}

	return status.GetState(), nil
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

	URL := os.Args[1]
	gitHubBaseURL := "https://github.com/"

	timer := time.Duration(60)
	if _, ok := os.LookupEnv("PRATUS_SLEEP_TIMER"); ok {
		customTimer, _ := strconv.Atoi(os.Getenv("PRATUS_SLEEP_TIMER"))
		timer = time.Duration(customTimer) * time.Second
	}

	sleepTimer := timer * time.Second

	owner, repo, number := parseGitHubURL(gitHubBaseURL, URL)

	fmt.Printf("Checking status of pull request %d in %s/%s every %d seconds\n", number, owner, repo, timer)

	for {

		state, err := getPRState(owner, repo, number)
		if err != nil {
			print(err)
		}

		if state != "pending" {
			fmt.Println("\nPR finished with state: " + state)
			os.Exit(0)
		}

		fmt.Print(".")
		time.Sleep(sleepTimer)
	}
}
