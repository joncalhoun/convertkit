package convertkit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joncalhoun/convertkit"
)

func TestClient_AuthError(t *testing.T) {
	check := func(t *testing.T, err error) {
		if err == nil {
			t.Fatalf("Account() err = nil; want error")
		}
		var ckErr convertkit.ErrorResponse
		if !errors.As(err, &ckErr) {
			t.Fatalf("Do() err type = %T; want %T", err, ckErr)
		}
		if ckErr.Type != "Authorization Failed" {
			t.Errorf("Type = %v; want %v", ckErr.Type, "Authorization Failed")
		}
		if ckErr.StatusCode != 401 {
			t.Errorf("StatusCode = %v; want %v", ckErr.StatusCode, 401)
		}
	}

	c := client(t, "fake-secret-key")
	c.Secret = "invalid"
	t.Run("GET", func(t *testing.T) {
		_, err := c.Account()
		check(t, err)
	})
	t.Run("DELETE", func(t *testing.T) {
		t.Skip("TODO: Implement this...")
	})
	t.Run("PUT", func(t *testing.T) {
		t.Skip("TODO: Implement this...")
	})
	t.Run("POST", func(t *testing.T) {
		t.Skip("TODO: Implement this...")
	})

}

func TestClient_AuthError_Post(t *testing.T) {
	c := client(t, "fake-secret-key")
	c.Secret = "invalid"
	_, err := c.Account()
	if err == nil {
		t.Fatalf("Account() err = nil; want error")
	}
	var ckErr convertkit.ErrorResponse
	if !errors.As(err, &ckErr) {
		t.Fatalf("Do() err type = %T; want %T", err, ckErr)
	}
	if ckErr.Type != "Authorization Failed" {
		t.Errorf("Type = %v; want %v", ckErr.Type, "Authorization Failed")
	}
	if ckErr.StatusCode != 401 {
		t.Errorf("StatusCode = %v; want %v", ckErr.StatusCode, 401)
	}
}

func client(t *testing.T, secret string) *convertkit.Client {
	c := clientWithHandler(t, baseHandler(t, secret))
	c.Secret = secret
	return c
}

func clientWithHandler(t *testing.T, handler http.HandlerFunc) *convertkit.Client {
	var c convertkit.Client
	server := httptest.NewServer(http.HandlerFunc(handler))
	c.BaseURL = server.URL
	t.Cleanup(func() {
		server.Close()
	})
	return &c
}

func baseHandler(t *testing.T, secret string) http.HandlerFunc {
	authErrHandler := testdataHandler(t, "ERR_auth")
	return func(w http.ResponseWriter, r *http.Request) {
		var apiSecret string
		switch r.Method {
		case http.MethodGet, http.MethodDelete:
			r.ParseForm()
			apiSecret = r.FormValue("api_secret")
		default:
			var temp struct {
				APISecret string `json:"api_secret"`
			}
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("reading body: %v", err)
			}
			r.Body.Close()
			br := bytes.NewReader(b)
			r.Body = ioutil.NopCloser(br)
			err = json.NewDecoder(br).Decode(&temp)
			if err != nil {
				t.Fatalf("Decode() err = %v; want nil", err)
			}
			apiSecret = temp.APISecret
		}
		if apiSecret != secret {
			t.Logf("Secret = %v; want %v", apiSecret, secret)
			authErrHandler(w, r)
			return
		}
		prefix := fmt.Sprintf("%s_%s", r.Method, strings.ReplaceAll(r.URL.Path[1:], "/", "_"))
		testdataHandler(t, prefix)(w, r)
	}
}

func testdataHandler(t *testing.T, prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode(t, prefix))
		body := body(t, prefix)
		defer body.Close()
		_, err := io.Copy(w, body)
		if err != nil {
			t.Fatalf("copying body to response: %v", err)
		}
	}
}

func statusCode(t *testing.T, path string) int {
	path = filepath.Join("testdata", fmt.Sprintf("%s.headers.json", path))
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return 200
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open the response body file: %s. err = %v", path, err)
	}
	var headers struct {
		StatusCode int `json:"status_code"`
	}
	err = json.NewDecoder(f).Decode(&headers)
	if err != nil {
		t.Fatalf("error decoding headers: %v", err)
	}
	return headers.StatusCode
}

func body(t *testing.T, path string) io.ReadCloser {
	path = filepath.Join("testdata", fmt.Sprintf("%s.json", path))
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open the response body file: %s. err = %v", path, err)
	}
	return f
}
