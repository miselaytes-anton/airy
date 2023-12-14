package testserver

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestServer struct {
	Server *httptest.Server
}

func (ts *TestServer) Get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Server.Client().Get(ts.Server.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

func (ts *TestServer) Post(t *testing.T, urlPath string, requestBody []byte) (int, http.Header, []byte) {
	r := bytes.NewReader(requestBody)
	rs, err := ts.Server.Client().Post(ts.Server.URL+urlPath, "application/json; charset=utf-8", r)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

func (ts *TestServer) Patch(t *testing.T, urlPath string, requestBody []byte) (int, http.Header, []byte) {
	r := bytes.NewReader(requestBody)

	req := httptest.NewRequest(
		http.MethodPatch,
		ts.Server.URL+urlPath,
		r,
	)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.RequestURI = ""

	rs, err := ts.Server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
