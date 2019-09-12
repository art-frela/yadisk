package yadisk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

const (
	requestTimeOut = 30 * time.Second
)

// YaDisk - main object for Yandex Disk
type YaDisk struct {
	url   string
	token string
	dlURL string
}

// YaLink - response format for get link YaDisk
type YaLink struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

// NewYaDisk - builder YaDisk
func NewYaDisk(url, token string) *YaDisk {
	ya := new(YaDisk)
	ya.url = url
	ya.token = token
	return ya
}

func (ya *YaDisk) getLink(path string) (*YaDisk, error) {
	uri := fmt.Sprintf("%s/?path=%s", ya.url, path)
	response, _, err := httpRequest(uri, ya.token)
	if err != nil {
		return ya, err
	}
	link := new(YaLink)
	err = json.NewDecoder(response).Decode(link)
	if err != nil {
		return ya, err
	}
	ya.dlURL = link.Href
	return ya, nil
}

// SaveToDisk - implement download file from YaDisk
func (ya *YaDisk) SaveToDisk(pathToFile string) error {
	_, err := ya.getLink(pathToFile)
	if err != nil {
		return err
	}
	dlReader, _, err := httpRequest(ya.dlURL, ya.token)
	newFN := path.Base(pathToFile)
	destFile, err := os.Create(newFN)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, dlReader)
	return err
}

// httpRequest - common part of http request
func httpRequest(uri, token string) (rc io.Reader, httpcode int, err error) {
	// set http client: timeout of request and switch off redirect
	c := http.Client{
		Timeout: requestTimeOut,
	}
	// make encoded url string, for cyrillic symbols and other
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	q := u.Query()
	u.RawQuery = q.Encode()
	// make request
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return
	}
	auth := fmt.Sprintf("OAuth %s", token)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "SampleGoClient/1.0")
	request.Header.Set("Authorization", auth)

	httpData, err := c.Do(request)
	if err != nil {
		return
	}
	httpcode = httpData.StatusCode
	if httpData.StatusCode != http.StatusOK {
		err = fmt.Errorf("some error at time http.request, request=%s; httpcode=%d", u.String(), httpData.StatusCode)
		return
	}
	rc = httpData.Body
	return
}
