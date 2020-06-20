package e2e

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/user/andon-webapp-in-go/src/db"
	"github.com/user/andon-webapp-in-go/src/view"
	"github.com/user/andon-webapp-in-go/src/wc"
)

func TestWorkcenterPage(t *testing.T) {
	database, err := db.GetDB()
	if err != nil {
		t.Fatalf("could not connect to database: %v", err)
	}
	wc.SetDB(database)
	view.RegisterStaticHandlers()
	h := wc.NewViewHandler()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse("http://localhost/wc/1")
		if err != nil {
			t.Fatalf("Failed to parse test URL: %v", err)
		}
		r.URL = url
		h.ServeHTTP(w, r)
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	defer resp.Body.Close()
	log.Println(string(body))
	if strings.Contains(string(body), "Assembly Line 1") == false {
		t.Error("Did not get expected response\n")
	}
}
