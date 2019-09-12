/*Package yadisk - simple implementation of Yandex Disk REST API
  Implemented methods:
	  download file
	  download public file

  Author: Artem Karpov, mailto: art.frela@gmail.com
  Date: 2019-09-12
  Subject: Geekbrains Go course, 2nd qrt "Essentials of Golang"

  Task: * Напишите функцию, которая получает на вход публичную ссылку на файл с «Яндекс.Диска» и сохраняет полученный файл на диск пользователя.
*/
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

// YaResource - properties of file at the Yandex disk
type YaResource struct {
	AntivirusStatus string      `json:"antivirus_status"`
	PublicKey       string      `json:"public_key"`
	PublicURL       string      `json:"public_url"`
	ViewsCount      int64       `json:"views_count"`
	Exif            interface{} `json:"exif"`
	Created         string      `json:"created"`
	Type            string      `json:"type"`
	ResourceID      string      `json:"resource_id"`
	Modified        string      `json:"modified"`
	Size            int64       `json:"size"`
	CommentIDS      interface{} `json:"comment_ids"`
	MIMEType        string      `json:"mime_type"`
	File            string      `json:"file"`
	MediaType       string      `json:"media_type"`
	Owner           interface{} `json:"owner"`
	Preview         string      `json:"preview"`
	Path            string      `json:"path"`
	Sha256          string      `json:"sha256"`
	Revision        int64       `json:"revision"`
	Md5             string      `json:"md5"`
	Name            string      `json:"name"`
}

// NewYaDisk - builder YaDisk
func NewYaDisk(uri, token string) *YaDisk {
	ya := new(YaDisk)
	ya.url = uri
	ya.token = token
	return ya
}

// NewPublicYaDisk - builder of public YaDisk
func NewPublicYaDisk(uri string) *YaDisk {
	ya := new(YaDisk)
	ya.url = uri
	return ya
}

// getDlLink - fill the download link attribute
func (ya *YaDisk) getDlLink(path string) error {
	uri := fmt.Sprintf("%s/disk/resources/download/?path=%s", ya.url, path)
	response, _, err := httpRequest(uri, ya.token)
	if err != nil {
		return err
	}
	link := new(YaLink)
	err = json.NewDecoder(response).Decode(link)
	if err != nil {
		return err
	}
	ya.dlURL = link.Href
	return nil
}

// getInfoPublicFile - fill download link and returns filename or error
func (ya *YaDisk) getInfoPublicFile(publicKey string) (filename string, err error) {
	uri := fmt.Sprintf("%s/disk/public/resources?public_key=%s", ya.url, publicKey)
	response, _, err := httpRequest(uri, ya.token)
	if err != nil {
		return
	}
	link := new(YaResource)
	err = json.NewDecoder(response).Decode(link)
	if err != nil {
		return
	}
	ya.dlURL = link.File
	filename = link.Name
	return
}

// SaveToDisk - implement download file from YaDisk, using token
func (ya *YaDisk) SaveToDisk(pathToFile string, isPublic bool) (err error) {
	if isPublic {
		pathToFile, err = ya.getInfoPublicFile(pathToFile)
	} else {
		err = ya.getDlLink(pathToFile)
	}
	if err != nil {
		return
	}
	dlReader, _, err := httpRequest(ya.dlURL, ya.token)
	if err != nil {
		return
	}
	newFN := path.Base(pathToFile)

	destFile, err := os.Create(newFN)
	if err != nil {
		return
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, dlReader)
	return
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
	if len(token) > 0 {
		request.Header.Set("Authorization", auth)
	}
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
