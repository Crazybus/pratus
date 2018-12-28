# Pratus

Waits for all statuses to have finished on a GitHub pull request. It can be combined with a notification tool like [noti](https://github.com/variadico/noti) to send you a notification once a pull request has finished testing.

## Usage

Set `GITHUB_TOKEN` to your personal GitHub token

```
$ export GITHUB_TOKEN=234oi2j3423j4io2342o34
```

Run pratus with the URL to the pull request

```
$ pratus https://github.com/Crazybus/pratus/pulls/1
Checking status of pull request 1 in Crazybus/pratus every 60 seconds
......
PR finished with state: success
```

## Configuration

* The sleep duration can be changed by setting `PRATUS_SLEEP_TIMER` to the amount of seconds (default `60`)
