package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hugozhu/godingtalk/demo/github"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/dnspod", Dnspod)
	http.HandleFunc("/github", github.Handle)
	// http.HandleFunc("/", handler)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func Dnspod(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("domain") + " " + r.FormValue("reason")
	token := r.FormValue("token")
	if token == "" {
		token = "b7e4b04c66b5d53669affb0b92cf533b9eff9b2bc47f86ff9f4227a2ba73798e"
	}
	request := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": msg,
		},
	}
	data, _ := json.Marshal(&request)
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token), bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", string(body))
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "http://hugozhu.myalert.info", 302)
// }