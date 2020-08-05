package util

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type RequestData struct {
	Query       map[string]string `json:"query"`
	BodyData    map[string]string `json:"body_data"`
	Header      map[string]string `json:"header"`
	ContentType string            `json:"content_type"`
}

// SendHTTPRequest send generic HTTP request, with url, method and request data.
func SendHTTPRequest(requestURL, method string, data *RequestData) ([]byte, error) {
	return nil, nil
}

// SendHTTPGetRequest send HTTP GET request.
func SendHTTPGetRequest(requestURL string) ([]byte, error) {
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// UnmarshalXMLBody unmarshal body data with XML.
func UnmarshalXMLBody(body []byte, data interface{}) error {
	return xml.Unmarshal(body, &data)
}
