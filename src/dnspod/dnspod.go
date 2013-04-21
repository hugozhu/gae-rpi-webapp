package dnspod

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Update(client *http.Client, cname string) string {
	body := url.Values{
		"login_email":    {login_email},
		"login_password": {login_password},
		"format":         {format},
		"domain_id":      {domain_id},
		"record_id":      {record_id},
		"sub_domain":     {sub_domain},
		"record_type":    {record_type},
		"record_line":    {record_line},
		"value":          {cname},
		"ttl":            {ttl},
	}
	req, err := http.NewRequest("POST", "https://dnsapi.cn/Record.Modify", strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "text/json")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	return string(bytes)
}
