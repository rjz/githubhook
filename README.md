githubhook
===============================================

Golang parser for [github webhooks][gh-webhook]. Not a server, though it could
be integrated with one.

[![Build Status](https://travis-ci.org/rjz/githubhook.svg?branch=master)](https://travis-ci.org/rjz/githubhook)
[![Go Report Card](https://goreportcard.com/badge/github.com/rjz/githubhook)](https://goreportcard.com/report/github.com/rjz/githubhook)
[![GoDoc](https://godoc.org/github.com/rjz/githubhook?status.svg)](https://godoc.org/github.com/rjz/githubhook)

Installation
-----------------------------------------------

```ShellSession
$ go get gopkg.in/rjz/githubhook.v0
```

Usage
-----------------------------------------------

Given an incoming `*http.Request` representing a webhook signed with a `secret`,
use `githubhook` to validate and parse its content:

```go
secret := []byte("don't tell!")
hook, err := githubhook.Parse(secret, req)
```

Plays nicely with the [google/go-github][gh-go-github] client!

```go
evt := github.PullRequestEvent{}
if err := json.Unmarshal(hook.Payload, &evt); err != nil {
  fmt.Println("Invalid JSON?", err)
}
```

[gh-webhook]: https://developer.github.com/webhooks/
[gh-go-github]: https://github.com/google/go-github
