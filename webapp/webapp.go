package webapp

import (
	"dnspod"
	"fmt"
	"io/ioutil"
	"net/http"
	"webapp/config"
	"webapp/counter"

	"appengine"
	"appengine/memcache"
	"appengine/urlfetch"
	"strconv"
)

func init() {
	http.HandleFunc(config.URL_BEACON, counter.Handle)

	http.HandleFunc("/online", counter.Count)
	http.HandleFunc("/online_get_token", counter.GetToken)
	http.HandleFunc("/online_send_msg", counter.SendMessage)

	http.HandleFunc("/switch_dns", switch_dns)

	http.HandleFunc("/ping", ping)

	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://hugozhu.myalert.info", 302)
}

func switch_dns(w http.ResponseWriter, r *http.Request) {
	cname := r.FormValue("cname")
	if cname == "" {
		cname = "gitcafe.myalert.info."
	}
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	w.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(w, "%s", dnspod.Update(client, cname))
}

func ping(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		url = "http://go.myalert.info/status.html"
	}
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	ok, body := url_ok(client, url)

	item, _ := memcache.Get(c, "site_fails")

	if ok {
		if item != nil {
			//switch back to pi
			dnspod.Update(client, "pi.myalert.info.")
			memcache.Delete(c, "site_fails")
		}
	} else {
		if item != nil {
			//previously failed, switch to github
			dnspod.Update(client, "gitcafe.myalert.info.")
			value, _ := strconv.Atoi(string(item.Value))
			value++
			item.Value = []byte(strconv.Itoa(value))
			memcache.Set(c, item)
		} else {
			//first time failed
			item = &memcache.Item{
				Key:   "site_fails",
				Value: []byte("0"),
			}
			memcache.Set(c, item)
		}
	}

	w.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(w, "%s", body)
}

func url_ok(client *http.Client, url string) (bool, string) {
	resp, err := client.Get(url)
	ok := false
	body := ""
	if err != nil {
		body = err.Error()
	} else {
		bytes, _ := ioutil.ReadAll(resp.Body)
		body = string(bytes)
		resp.Body.Close()
		ok = resp.StatusCode == 200
	}
	return ok, body
}
