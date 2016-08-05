package githubhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// Hook describes an inbound github webhook
type Hook struct {
	Signature string
	Event     string
	Id        string
	Payload   []byte
}

const signaturePrefix = "sha1="
const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

// SignedBy checks that the provided secret matches the hook Signature
func (h *Hook) SignedBy(secret []byte) bool {
	if len(h.Signature) != signatureLength || !strings.HasPrefix(h.Signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(h.Signature[5:]))

	return hmac.Equal(signBody(secret, h.Payload), actual)
}

// New extracts a Hook from an incoming http.Request
func New(req *http.Request) (hook *Hook, err error) {
	hook = new(Hook)
	if !strings.EqualFold(req.Method, "POST") {
		return nil, errors.New("Unknown method!")
	}

	if hook.Signature = req.Header.Get("x-hub-signature"); len(hook.Signature) == 0 {
		return nil, errors.New("No signature!")
	}

	if hook.Event = req.Header.Get("x-github-event"); len(hook.Event) == 0 {
		return nil, errors.New("No event!")
	}

	if hook.Id = req.Header.Get("x-github-delivery"); len(hook.Id) == 0 {
		return nil, errors.New("No event Id!")
	}

	hook.Payload, err = ioutil.ReadAll(req.Body)
	return
}

// Parse extracts and verifies a hook against a secret
func Parse(secret []byte, req *http.Request) (hook *Hook, err error) {
	hook, err = New(req)
	if err == nil && !hook.SignedBy(secret) {
		err = errors.New("Invalid signature")
	}
	return
}
