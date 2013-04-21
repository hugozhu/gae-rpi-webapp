package webapp

import (
	"dnspod"
	"fmt"
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

	http.HandleFunc("/switch_dns", switch_dns)

	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://hugozhu.myalert.info", 302)
}

func switch_dns(w http.ResponseWriter, r *http.Request) {
	cname := r.FormValue("cname")
	if cname == "" {
		cname = "hugozhu.github.io."
	}
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	w.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(w, "%s", dnspod.Update(client, cname))
}
