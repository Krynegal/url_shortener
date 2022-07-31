package handlers

import (
	"bytes"
	"fmt"
	"github.com/Krynegal/url_shortener.git/internal/configs"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	cfg    = configs.NewConfig()
	s, err = storage.NewStorage(cfg)
)

const (
	url1 = "http://someHost.ya/qwerty"
	url2 = "http://newHost.com/okmijn"
	url3 = "http:/hhhost.com/wsxedc"
)

func testRequest(t *testing.T, method, path string, body string) (*http.Response, string) {
	ts := httptest.NewServer(NewHandler(s, cfg).Mux)
	defer ts.Close()

	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer([]byte(body)))
	require.NoError(t, err)

	client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, strings.Trim(string(respBody), "\n")
}

func TestRouterPostPositive1(t *testing.T) {
	resp, body := testRequest(t, http.MethodPost, "/", url1)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/1", body)
}

func TestRouterPostPositive2(t *testing.T) {
	resp, body := testRequest(t, http.MethodPost, "/", url2)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/2", body)

}

func TestRouterPostNegative1(t *testing.T) {
	resp, body := testRequest(t, http.MethodPost, "/1", "body")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Equal(t, "", body)
}

func TestRouterGetPositive1(t *testing.T) {
	resp, _ := testRequest(t, http.MethodGet, "/1", "")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, url1, resp.Header.Get("Location"))
}

func TestRouterGetPositive2(t *testing.T) {
	resp, _ := testRequest(t, http.MethodGet, "/2", "")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, url2, resp.Header.Get("Location"))
}

func TestRouterGetNegative1(t *testing.T) {
	resp, _ := testRequest(t, http.MethodPost, "/2", "")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Equal(t, "", resp.Header.Get("Location"))
}

func TestRouterPostJSONpositive1(t *testing.T) {
	reqBody := fmt.Sprintf(`{"url":"%s"}`, url3)
	resp, body := testRequest(t, http.MethodPost, "/api/shorten", reqBody)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, `{"result":"http://localhost:8080/3"}`, body)
}

func TestRouterPostJSONnegative1(t *testing.T) {
	resp, body := testRequest(t, http.MethodPost, "/api/shorten", "")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "wrong request format", body)
}

func TestRouterPostJSONnegative2(t *testing.T) {
	reqBody := fmt.Sprintf(`{"url":"%s"}`, "")
	resp, body := testRequest(t, http.MethodPost, "/api/shorten", reqBody)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "url is empty", body)
}
