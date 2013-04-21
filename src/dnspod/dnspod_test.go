package dnspod

import (
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
	client := &http.Client{}
	t.Log(Update("hugozhu.github.com."))
	t.Log(Update("pi.myalert.info."))
}
