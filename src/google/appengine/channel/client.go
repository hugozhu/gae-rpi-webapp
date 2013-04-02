package channel

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func NewChannel(key string, token_url string) (c *Channel) {
	c = &Channel{
		Key:           key,
		URL_Get_Token: token_url,
		User_Agent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.43 Safari/537.31",
		Params: map[string]string{
			"host": "124", //todo: change to random host
		},
		Client: &http.Client{},
		Scookie: &http.Cookie{
			Name:  "S",
			Value: "",
		},
		rid: 1580,
	}
	return
}

type DefaultChannelHandler struct {
}

func (DefaultChannelHandler) OnOpened() {

}
func (DefaultChannelHandler) OnMessage(msg string) {

}
func (DefaultChannelHandler) onError(err error) {

}
func (DefaultChannelHandler) onClose() {

}

func (c *Channel) Open() ChannelSocket {
	c.Params["token"] = "AHRlWrqFMYxLsivdKDcDuWL7vrus4lE_gBI0tQYIPuedVOyhvhJTZTqzeG8iMq4Ks3LdP5p5wtwzfAM_1u4RBd9l4-8cuT9O-Q"
	// c.Params["token"] = c.NewToken()

	c.get_clid_gsessionid()
	c.test_clid_gsessionid()
	c.get_sid()                 //get SID
	c.register_new_conneciton() //register a new connection
	c.receive()

	return DefaultChannelHandler{}
}

func (c *Channel) NewToken() string {
	url := c.URL_Get_Token + c.Key
	resp := c.HttpCall(url, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	token := strings.TrimSpace(string(body))
	log.Println("new token", token)
	c.Params["token"] = token
	return token
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

func (c *Channel) HttpCall(_url string, body io.Reader) (resp *http.Response) {
	c.Params["zx"] = RandomString(12)
	_url = make_url(_url, c.Params)
	method := "GET"
	if body != nil {
		method = "POST"
	}
	req, err := http.NewRequest(method, _url, body)
	req.Header.Set("User-Agent", c.User_Agent)
	req.Header.Set("Referer", c.URL_Get_Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if err != nil {
		panic(err)
	}
	if c.Scookie.Value != "" {
		req.AddCookie(c.Scookie)
	}
	resp, err = c.Client.Do(req)
	if err != nil {
		panic(err)
	}
	return
}

func (c *Channel) get_clid_gsessionid() {
	resp := c.HttpCall("http://${host}.talkgadget.google.com/talkgadget/d?token=${token}&xpc=%7B%22cn%22%3A%22kb7TjvGhBn%22%2C%22tp%22%3Anull%2C%22osh%22%3Anull%2C%22ppu%22%3A%22http%3A%2F%2Fapp.myalert.info%2F_ah%2Fchannel%2Fxpc_blank%22%2C%22lpu%22%3A%22http%3A%2F%2F${host}.talkgadget.google.com%2Ftalkgadget%2Fxpc_blank%22%7D", nil)
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if "S" == cookie.Name {
			c.Scookie.Value = cookie.Value
			break
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	s := string(body)
	arr := strings.Split(s[strings.LastIndex(s, "javascript"):], "\n")
	c.Params["clid"] = strings.Trim(arr[3], "\", \t\r\n")
	c.Params["gsessionid"] = strings.Trim(arr[4], "\", \t\r\n")
}

func (c *Channel) test_clid_gsessionid() {
	resp := c.HttpCall("http://${host}.talkgadget.google.com/talkgadget/dch/test?VER=8&clid=${clid}&gsessionid=${gsessionid}&prop=data&token=${token}&ec=%5B%22ci%3Aec%22%5D&MODE=init&zx=${zx}&t=1", nil)
	defer resp.Body.Close()
}

var sidRe = regexp.MustCompile(`(?ismU),\"(.*)\"`)

func (c *Channel) get_sid() {
	c.Params["RID"] = fmt.Sprintf("%d", c.rid)
	postfields := &url.Values{
		"count": []string{"0"},
	}
	resp := c.HttpCall("http://${host}.talkgadget.google.com/talkgadget/dch/bind?VER=8&clid=${clid}&gsessionid=${gsessionid}&prop=data&token=${token}&ec=%5B%22ci%3Aec%22%5D&RID=${RID}&CVER=1&zx=${zx}&t=1", strings.NewReader(postfields.Encode()))

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	s := string(body)
	lines := strings.Split(s, "\n")
	matches := sidRe.FindStringSubmatch(lines[1])
	if len(matches) == 0 {
		panic("Token expired")
	}
	if matches[1] != "" {
		c.Params["SID"] = matches[1]
	} else {
		panic("Can't get SID")
	}
	log.Println("SID", c.Params["SID"])
}

func (c *Channel) register_new_conneciton() {
	c.rid++
	c.Params["RID"] = fmt.Sprintf("%d", c.rid)

	postfields := &url.Values{}
	postfields.Set("count", "1")
	postfields.Set("ofs", "0")
	postfields.Set("req0_m", `["connect-add-client"]`)
	postfields.Set("req0_c", c.Params["clid"])
	postfields.Set("req0__sc", "c")

	resp := c.HttpCall("http://${host}.talkgadget.google.com/talkgadget/dch/bind?VER=8&clid=${clid}&gsessionid=${gsessionid}&prop=data&token=${token}&ec=%5B%22ci%3Aec%22%5D&SID=${SID}&RID=${RID}&AID=2&zx=${zx}&t=1", strings.NewReader(postfields.Encode()))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	c.Handler.OnOpened()
	log.Println("register_new_conneciton response:\n", string(body))
}

func (c *Channel) receive() {
	resp := c.HttpCall("http://${host}.talkgadget.google.com/talkgadget/dch/bind?VER=8&clid=${clid}&gsessionid=${gsessionid}&prop=data&token=${token}&ec=%5B%22ci%3Aec%22%5D&RID=rpc&SID=${SID}&CI=0&AID=2&TYPE=xmlhttp&zx=${zx}&t=1", nil)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {

			} else {
				c.Handler.onError(err)
			}
		}
		len, err := strconv.Atoi(string(bytes))
		if err != nil {
			panic(err)
		}
		buf := make([]byte, len)
		n, err := reader.Read(buf)
		if err != nil {
			panic(err)
		}
		if n == len {
			c.Handler.OnMessage(string(buf))
		} else {
			panic("Read less ?")
		}
	}
	c.Handler.onClose()
}
