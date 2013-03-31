package webapp

import (
	"net/http"
	"webapp/config"
	"webapp/counter"
)

func init() {
	http.HandleFunc(config.URL_BEACON, counter.Handle)
	http.HandleFunc("/online", counter.Count)
	http.HandleFunc("/online_get_token", counter.GetToken)
	http.HandleFunc("/online_send_msg", counter.SendMessage)

	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://hugozhu.myalert.info", 302)
}
