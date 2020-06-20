package hash_test

import (
	"testing"

	"github.com/user/andon-webapp-in-go/src/hash"
)

func TestHash(t *testing.T) {
	data := []struct {
		have []string
		want string
	}{
		{have: []string{"foo"}, want: "lTPclwvmB9CX/QeuBlKZNvsiZLtB5i9PCE7fIMOKLVlkbYFuraSgPIHyzaHZKL432sWAYECTYQpzl2zuKYGHQA=="},
		{have: []string{"foo", "bar"}, want: "9Wyi4V42A8h4hceRu6sBYySzjNE4pYOkRRjm7EfN/F3UkLQsyjS/MJV/qAcWECf1BQyMYcuTHvdCL6AQGwlBeA=="},
	}

	for _, item := range data {
		got := hash.Hash(item.have...)
		if got != item.want {
			t.Errorf("Unexpected hash result for %v. Got %q, wanted %q", item.have, got, item.want)
		}
	}
}
