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

[gh-webhook]: https://developer.github.com/webhooks/
