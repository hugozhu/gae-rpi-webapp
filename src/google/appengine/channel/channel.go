package channel

import (
	"net/http"
)

type Channel struct {
	Key           string
	URL_Get_Token string
	User_Agent    string
	Params        map[string]string
	Client        *http.Client
	Scookie       *http.Cookie
	Handler       *ChannelSocket

	rid int
}

type ChannelSocket struct {
	OnOpened  func()
	OnMessage func(msg string)
	OnError   func(err error)
	OnClose   func()
}
