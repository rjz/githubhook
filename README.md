githubhook
===============================================

Golang parser for [github webhooks][gh-webhook]. Not a server, though it could
be integrated with one.

Installation
-----------------------------------------------

    $ go get github.com/rjz/githubhook

Usage
-----------------------------------------------

Given an incoming `*http.Request` representing a webhook signed with a `secret`,
use `githubhook` to validate and parse its content:

    secret := []byte("don't tell!")
    hook, err := githubhook.Parse(secret, req)

Plays nicely with the [google/go-github][gh-go-github] client!

    evt := github.PullRequestEvent{}
    if err := json.Unmarshal(hook.Payload, &evt); err != nil {
      fmt.Println("Invalid JSON?", err)
    }

[gh-webhook]: https://developer.github.com/webhooks/
[gh-go-github]: https://github.com/google/go-github
