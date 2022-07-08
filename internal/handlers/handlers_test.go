package handlers

import (
	"bytes"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var s = storage.NewStorage(storage.Memory)

func testRequest(t *testing.T, method, path string, body string) (*http.Response, string) {
	ts := httptest.NewServer(NewHandler(s).Mux)
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

	return resp, string(respBody)
}

func TestRouterPOST(t *testing.T) {
	resp, body := testRequest(t, http.MethodPost, "/", "http://someHost")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/1", body)

}

func TestRouterPOST2(t *testing.T) {
	resp, body := testRequest(t, http.MethodPost, "/", "http://newHost.com/okmijn")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/2", body)

}

func TestRouterGet(t *testing.T) {
	resp, body := testRequest(t, http.MethodGet, "/1", "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "http://someHost", resp.Header.Get("Location"))
	assert.Equal(t, "", body)
}

func TestRouterGet2(t *testing.T) {
	resp, body := testRequest(t, http.MethodGet, "/2", "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "http://newHost.com/okmijn", resp.Header.Get("Location"))
	assert.Equal(t, "", body)
}

func TestRouterPOST1(t *testing.T) {
	resp, body := testRequest(t, http.MethodPost, "/1", "body")
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Equal(t, "", body)
}

//func TestStatusHandler(t *testing.T) {
//	s := storage.NewStorage(storage.Memory)
//
//	// определяем структуру теста
//	type want struct {
//		code     int
//		response string
//	}
//	// создаём массив тестов: имя и желаемый результат
//	tests := []struct {
//		name    string
//		handler http.HandlerFunc
//		method  string
//		target  string
//		body    string
//		want    want
//	}{
//		// определяем все тесты
//		{
//			name:    "positive test for POST request",
//			handler: ShortURL(s),
//			method:  http.MethodGet,
//			target:  "http://localhost:8080/",
//			body:    "http://someHost.ya/qwerty",
//			want: want{
//				code:     201,
//				response: "http://localhost:8080/1",
//			},
//		},
//		{
//			name:    "positive test for GET request",
//			handler: GetID(s),
//			method:  http.MethodGet,
//			target:  "http://localhost:8080/1",
//			body:    "",
//			want: want{
//				code:     307,
//				response: "",
//			},
//		},
//	}
//	for _, tt := range tests {
//		// запускаем каждый тест
//		t.Run(tt.name, func(t *testing.T) {
//			router := mux.NewRouter()
//			ts := httptest.NewServer(router)
//			request, err := http.NewRequest(tt.method, ts.URL+tt.target, bytes.NewBuffer([]byte(tt.body)))
//			if err != nil {
//				panic("panic 1")
//			}
//			res, err := http.DefaultClient.Do(request)
//			// создаём новый Recorder
//			//w := httptest.NewRecorder()
//			// определяем хендлер
//			//h := http.Handler(tt.handler)
//			// запускаем сервер
//			//h.ServeHTTP(w, request)
//			//res := w.Result()
//
//			log.Printf("res: %v", res)
//			// проверяем код ответа
//			assert.Equal(t, tt.want.code, res.StatusCode)
//
//			// получаем и проверяем тело запроса
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//			if err != nil {
//				t.Fatal(err)
//			}
//			if string(resBody) != tt.want.response {
//				t.Errorf("Expected body %s, got %s", tt.want.response, string(resBody))
//			}
//		})
//	}
//}
