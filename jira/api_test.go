package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubHTTPHandler struct {
	Response string
}

func (stub StubHTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, stub.Response)
}

func TestHttpApiReturnsBodyContents(t *testing.T) {
	response := "{message: 'Hello, World'}"
	stubHandler := StubHTTPHandler{Response: response}
	server := httptest.NewServer(stubHandler)
	defer server.Close()
	api := new(HttpApi)
	output, err := api.Fetch(server.URL, map[string]string{})
	outputString := fmt.Sprintf("%s", output)
	responseString := fmt.Sprintf("%s\n", response)
	if err != nil {
		t.Error("Expected API to not produce error", err.Error())
	}
	if outputString != responseString {
		t.Error("Expected api to forward response through", outputString)
	}
}
