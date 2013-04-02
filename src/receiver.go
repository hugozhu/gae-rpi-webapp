package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	channel = "pi"
	params  = map[string]string{
		"host": "124",
	}

	client = &http.Client{}

	scookie = &http.Cookie{
		Name:  "S",
		Value: "",
	}
	RID = 1580
)

func init() {
	log.Println("starting ....")
	params["token"] = "AHRlWrqFMYxLsivdKDcDuWL7vrus4lE_gBI0tQYIPuedVOyhvhJTZTqzeG8iMq4Ks3LdP5p5wtwzfAM_1u4RBd9l4-8cuT9O-Q"
	//params["token"] = get_token(channel)
}

func main() {
	get_clid_gsessionid()
	test_clid_gsessionid()
	get_sid()                 //get SID
	register_new_conneciton() //register a new connection
	receive()
}

func make_url(url string, params map[string]string) string {
	for k, v := range params {
		url = strings.Replace(url, "${"+k+"}", v, -1)
	}
	return url
}

func RandomString(len int) string {
	return "h3o58yoptmcc"
}

func get_token(channel string) string {
	url := "http://app.myalert.info/online_get_token?id=" + channel
	resp := HttpCall(url, nil)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	token := strings.TrimSpace(string(body))
	log.Println("token", token)
	return token
}

func get_clid_gsessionid() {
	resp := HttpCall("http://${host}.talkgadget.google.com/talkgadget/d?token=${token}&xpc=%7B%22cn%22%3A%22kb7TjvGhBn%22%2C%22tp%22%3Anull%2C%22osh%22%3Anull%2C%22ppu%22%3A%22http%3A%2F%2Fapp.myalert.info%2F_ah%2Fchannel%2Fxpc_blank%22%2C%22lpu%22%3A%22http%3A%2F%2F${host}.talkgadget.google.com%2Ftalkgadget%2Fxpc_blank%22%7D", nil)
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if "S" == c.Name {
			scookie.Value = c.Value
			break
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	s := string(body)
	arr := strings.Split(s[strings.LastIndex(s, "javascript"):], "\n")
	params["clid"] = strings.Trim(arr[3], "\", \t\r\n")
	params["gsessionid"] = strings.Trim(arr[4], "\", \t\r\n")
}

func test_clid_gsessionid() {
	resp := HttpCall("http://${host}.talkgadget.google.com/talkgadget/dch/test?VER=8&clid=${clid}&gsessionid=${gsessionid}&prop=data&token=${token}&ec=%5B%22ci%3Aec%22%5D&MODE=init&zx=${zx}&t=1", nil)
	defer resp.Body.Close()
}

var sidRe = regexp.MustCompile(`(?ismU),\"(.*)\"`)

func get_sid() {
	params["RID"] = fmt.Sprintf("%d", RID)
	postfields := &url.Values{
		"count": []string{"0"},
	}
	resp := HttpCall("http://${host}.talkgadget.google.com/talkgadget/dch/bind?VER=8&clid=${clid}&gsessionid=${gsessionid}&prop=data&token=${token}&ec=%5B%22ci%3Aec%22%5D&RID=${RID}&CVER=1&zx=${zx}&t=1", strings.NewReader(postfields.Encode()))

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	s := string(body)
	lines := strings.Split(s, "\n")
	matches := sidRe.FindStringSubmatch(lines[1])
	if matches[1] != "" {
		params["SID"] = matches[1]
	} else {
		panic("Can't get SID")
	}
	log.Println("SID", params["SID"])
}

func register_new_conneciton() {
	RID++
	params["RID"] = fmt.Sprintf("%d", RID)

	postfields := &url.Values{}
	postfields.Set("count", "1")
	postfields.Set("ofs", "0")
	postfields.Set("req0_m", `["connect-add-client"]`)
	postfields.Set("req0_c", params["clid"])
	postfields.Set("req0__sc", "c")

	resp := HttpCall("http://${host}.talkgadget.google.com/talkgadget/dch/bind?VER=8&clid=${clid}&gsessionid=${gsessionid}&prop=data&token=${token}&ec=%5B%22ci%3Aec%22%5D&SID=${SID}&RID=${RID}&AID=2&zx=${zx}&t=1", strings.NewReader(postfields.Encode()))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("register_new_conneciton response:\n", string(body))
}

func receive() {
	resp := HttpCall("http://${host}.talkgadget.google.com/talkgadget/dch/bind?VER=8&clid=${clid}&gsessionid=${gsessionid}&prop=data&token=${token}&ec=%5B%22ci%3Aec%22%5D&RID=rpc&SID=${SID}&CI=0&AID=2&TYPE=xmlhttp&zx=${zx}&t=1", nil)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {

			} else {
				panic(err)
			}
		}
		log.Println(string(bytes))
	}
}

func HttpCall(_url string, body io.Reader) (resp *http.Response) {
	params["zx"] = RandomString(12)
	_url = make_url(_url, params)
	method := "GET"
	if body != nil {
		method = "POST"
	}
	req, err := http.NewRequest(method, _url, body)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.43 Safari/537.31")
	req.Header.Set("Referer", "http://app.myalert.info/online.html")
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		panic(err)
	}
	if scookie.Value != "" {
		req.AddCookie(scookie)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	return
}
