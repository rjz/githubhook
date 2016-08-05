package githubhook_test

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/rjz/githubhook"
	"net/http"
	"strings"
	"testing"
)

const testSecret = "foobar"

func expectErrorMessage(t *testing.T, msg string, err error) {
	if err == nil || err.Error() != msg {
		t.Error(fmt.Sprintf("Expected '%s', got %s", msg, err))
	}
}

func expectNewError(t *testing.T, msg string, r *http.Request) {
	_, err := githubhook.New(r)
	expectErrorMessage(t, msg, err)
}

func expectParseError(t *testing.T, msg string, r *http.Request) {
	_, err := githubhook.Parse([]byte(testSecret), r)
	expectErrorMessage(t, msg, err)
}

func signature(body string) string {
	dst := make([]byte, 40)
	computed := hmac.New(sha1.New, []byte(testSecret))
	computed.Write([]byte(body))
	hex.Encode(dst, computed.Sum(nil))
	return "sha1=" + string(dst)
}

func TestNonPost(t *testing.T) {
	r, _ := http.NewRequest("GET", "/path", nil)
	expectNewError(t, "Unknown method!", r)
}

func TestMissingSignature(t *testing.T) {
	r, _ := http.NewRequest("POST", "/path", nil)
	expectNewError(t, "No signature!", r)
}

func TestMissingEvent(t *testing.T) {
	r, _ := http.NewRequest("POST", "/path", nil)
	r.Header.Add("x-hub-signature", "bogus signature")
	expectNewError(t, "No event!", r)
}

func TestMissingEventId(t *testing.T) {
	r, _ := http.NewRequest("POST", "/path", nil)
	r.Header.Add("x-hub-signature", "bogus signature")
	r.Header.Add("x-github-event", "bogus event")
	expectNewError(t, "No event Id!", r)
}

func TestInvalidSignature(t *testing.T) {
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("..."))
	r.Header.Add("x-hub-signature", "bogus signature")
	r.Header.Add("x-github-event", "bogus event")
	r.Header.Add("x-github-delivery", "bogus id")
	expectParseError(t, "Invalid signature", r)
}

func TestValidSignature(t *testing.T) {

	body := "{}"

	r, _ := http.NewRequest("POST", "/path", strings.NewReader(body))
	r.Header.Add("x-hub-signature", signature(body))
	r.Header.Add("x-github-event", "bogus event")
	r.Header.Add("x-github-delivery", "bogus id")

	if _, err := githubhook.Parse([]byte(testSecret), r); err != nil {
		t.Error(fmt.Sprintf("Unexpected error '%s'", err))
	}
}
