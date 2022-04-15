package util

import (
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// SendHTTPGetRequest send HTTP GET request.
func SendHTTPGetRequest(requestURL string) ([]byte, error) {
	// resp, err := http.Get(requestURL)
	tr := &http.Transport{ // solve x509: certificate(https)
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(requestURL)
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
	return xml.Unmarshal(body, data)
}
