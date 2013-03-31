package counter

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/channel"
	"appengine/log"

	"webapp/config"
)

var CHARS = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
	"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
	"u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y", "Z"}

var BASE = int64(len(CHARS))

func str2int(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func make_hash(serverip string, ip string, timestamp int64) string {
	ip1 := IP2Long(serverip)
	ip2 := IP2Long(ip)
	crc := (ip1 + ip2 + timestamp) % config.SALT

	return fmt.Sprintf("%s%s%s%s", Long2String(ip1), Long2String(ip2), Long2String(timestamp), Long2String(crc))
}

func IP2Long(ip string) int64 {
	parts := strings.SplitN(ip, ".", 4)
	if len(parts) != 4 {
		return 0
	}
	return str2int(parts[0])<<24 + str2int(parts[1])<<16 + str2int(parts[2])<<8 + str2int(parts[3])
}

func Long2IP(i int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", i>>24, i&0x00FFFFFF>>16, i&0x0000FFFF>>8, i&0x000000FF)
}

func Long2String(id int64) string {
	if id < BASE {
		return CHARS[id]
	}
	n := id % BASE
	return Long2String(id/BASE) + CHARS[n]
}

func String2Long(s string) int64 {
	//not implemented yet
	return 0
}

func localIP() (net.IP, error) {
	tt, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range aa {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if v4 == nil || v4[0] == 127 { // loopback address
				continue
			}
			return v4, nil
		}
	}
	return nil, errors.New("cannot find local IP address")
}

var GIF []byte

func init() {
	GIF, _ = base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
}

func Handle(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)

	now := time.Now()
	expire := now.AddDate(30, 0, 0)
	zcookie, _ := r.Cookie("z")
	if zcookie == nil {
		zcookie = &http.Cookie{}
		zcookie.Name = "z"
		zcookie.Value = make_hash("127.0.0.1", r.RemoteAddr, now.UnixNano())
		zcookie.Expires = expire
		zcookie.Path = "/"
		zcookie.Domain = config.DOMAIN
		http.SetCookie(w, zcookie)
	}
	context.Infof("%s", zcookie.Value)

	w.Header().Set("Content-type", "image/gif")
	w.Header().Set("Cache-control", "no-cache, must-revalidate")
	w.Header().Set("Expires", "Sat, 26 Jul 1997 05:00:00 GMT")

	fmt.Fprintf(w, "%s", GIF)
}

func Count(w http.ResponseWriter, r *http.Request) {
	count := 0
	uniq := make(map[string]bool)
	c := appengine.NewContext(r)
	duration := 5
	if r.FormValue("duration") != "" {
		duration, _ = strconv.Atoi(r.FormValue("duration"))
	}
	query := &log.Query{
		AppLogs:   true,
		StartTime: time.Now().Add(time.Duration(-1*duration) * time.Minute),
		Versions:  []string{"1"},
	}

	for results := query.Run(c); ; {
		record, err := results.Next()
		if err == log.Done {
			break
		}
		if err != nil {
			c.Errorf("Failed to retrieve next log: %v", err)
		}

		if len(record.AppLogs) > 0 && strings.Index(record.Combined, "GET "+config.URL_BEACON) > 0 {
			zcookie := record.AppLogs[0].Message
			if zcookie != "" {
				count++
				uniq[zcookie] = true
			}
		}
	}
	fmt.Fprintf(w, "%d %d\n", len(uniq), count)
}

var tpl = `
<html>
    <body>
        <script type="text/javascript" src="/_ah/channel/jsapi"></script>

        <script>
            function onOpened() {
                alert("onOpened")
            }

            function onMessage(x) {
                alert("onMessage "+x)
            }

            function onError(x) {
                alert("onError "+x)
            }

            function onClose(x) {
                alert("onClose "+x)
            }

            channel = new goog.appengine.Channel('{{.token}}');
            socket = channel.open();
            socket.onopen = onOpened;
            socket.onmessage = onMessage;
            socket.onerror = onError;
            socket.onclose = onClose;
        </script>
    </body>
</html>
`
var mainTemplate = template.Must(template.New("main").Parse(tpl))

func GetToken(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	tok, err := channel.Create(c, r.FormValue("id"))
	if err != nil {
		http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
		c.Errorf("channel.Create: %v", err)
		return
	}
	// fmt.Fprintf(w, "%s", tok)

	mainTemplate.Execute(w, map[string]string{
		"token": tok,
	})
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	token := r.FormValue("id")
	channel.SendJSON(c, token, "{1,3,4}")
}
