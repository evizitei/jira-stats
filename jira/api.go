package jira

import (
	"io/ioutil"
	"net/http"
)

type Api interface {
	Fetch(url string, headers map[string]string) ([]byte, error)
}

type HttpApi struct {
}

func (HttpApi) Fetch(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println("Request Building Error:", err.Error())
		return nil, err
	}

	for header, value := range headers {
		request.Header.Add(header, value)
	}
	response, err := client.Do(request)
	if err != nil {
		println("API Connection Error:", err.Error())
		return nil, err
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		println("Response Reading Error:", err.Error())
		return nil, err
	}

	return contents, nil
}
