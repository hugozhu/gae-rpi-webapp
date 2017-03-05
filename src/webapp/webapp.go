package webapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"webapp/config"
	"webapp/counter"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc(config.URL_BEACON, counter.Handle)
	http.HandleFunc("/online", counter.Count)
	http.HandleFunc("/online_get_token", counter.GetToken)
	http.HandleFunc("/online_send_msg", counter.SendMessage)
	http.HandleFunc("/dnspod", Dnspod)
	// http.HandleFunc("/", handler)
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
