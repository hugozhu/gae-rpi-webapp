package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"appengine"
	"appengine/urlfetch"

	"strings"

	"github.com/google/go-github/github"
	dingtalk "github.com/hugozhu/godingtalk"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

const CHAT_ID = "chat6a93bc1ee3b7d660d372b1b877a9de62"
const SENDER_ID = "011217462940"

func Handle(w http.ResponseWriter, r *http.Request) {
	var err error

	action := r.Header.Get("x-github-event")
	signature := r.Header.Get("x-hub-signature")
	// id := r.Header.Get("x-github-delivery")
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if !verifySignature([]byte(os.Getenv("GITHUB_WEBHOOK_SECRET")), signature, body) {
		http.Error(w, fmt.Sprintf("%v", "Invalid or empty signature"), http.StatusBadRequest)
		return
	}

	var event github.WebHookPayload
	json.Unmarshal(body, &event)

	context := appengine.NewContext(r)
	c := dingtalk.NewDingTalkClient(os.Getenv("CORP_ID"), os.Getenv("CORP_SECRET"))
	c.HTTPClient = urlfetch.Client(context)
	c.HTTPClient.Transport = &urlfetch.Transport{
		Context: context,
		AllowInvalidServerCertificate: true,
	}
	c.RefreshAccessToken()

	msg := dingtalk.OAMessage{}
	msg.Head.Text = "Github"
	msg.Head.BgColor = "FF00AABB"
	switch action {
	case "push":
		msg.Body.Title = "[" + *event.Repo.Name + "] Push"
		msg.URL = *event.Compare
		for _, commit := range event.Commits {
			msg.Body.Form = append(msg.Body.Form, dingtalk.OAMessageForm{
				Key:   "Commits: ",
				Value: *commit.Message + "\n Modified: " + strings.Join(commit.Modified, ","),
			})
		}
	case "watch":
		msg.Body.Title = "[" + *event.Repo.Name + "] Watch Updated"
		msg.URL = *event.Repo.HTMLURL
		msg.Body.Form = []dingtalk.OAMessageForm{
			{
				Key:   "Watchers: ",
				Value: fmt.Sprintf("%d", *event.Repo.WatchersCount),
			},
			{
				Key:   "Forks: ",
				Value: fmt.Sprintf("%d", *event.Repo.ForksCount),
			},
		}
	default:
		msg.Body.Title = "[" + *event.Repo.Name + "] " + action
		msg.URL = *event.Repo.HTMLURL
	}
	msg.Body.Author = *event.Sender.Login

	err = c.SendOAMessage(SENDER_ID, CHAT_ID, msg)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}
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
