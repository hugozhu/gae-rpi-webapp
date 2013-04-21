package dnspod

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	Update("hugozhu.github.com.")
	Update("pi.myalert.info.")
}
