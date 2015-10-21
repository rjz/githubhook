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

type Hook struct {
	Signature string
	Event     string
	Id        string
	Payload   []byte
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

func Parse(secret []byte, req *http.Request) (*Hook, error) {
	hook := Hook{}

	if hook.Signature = req.Header.Get("x-hub-signature"); len(hook.Signature) == 0 {
		return nil, errors.New("No signature!")
	}

	if hook.Event = req.Header.Get("x-github-event"); len(hook.Event) == 0 {
		return nil, errors.New("No event!")
	}

	if hook.Id = req.Header.Get("x-github-delivery"); len(hook.Id) == 0 {
		return nil, errors.New("No event Id!")
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	if !verifySignature(secret, hook.Signature, body) {
		return nil, errors.New("Invalid signature")
	}

	hook.Payload = body

	return &hook, nil
}
