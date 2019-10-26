package http

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 100,
	},
	Timeout: 10 * time.Second,
}

func Send(req *http.Request) ([]byte, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return body, err
}

func Get(url string, header http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	return Send(req)
}

func Post(url string, body []byte, header http.Header) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header = header
	return Send(req)
}

func PostWithJson(url string, body []byte, header http.Header) ([]byte, error) {
	header.Set("Content-Type", "application/json")
	return Post(url, body, header)
}

func WrapSignHeader(body []byte, key, secret string) http.Header {
	header := make(http.Header)
	ts := fmt.Sprintf("%d", time.Now().Unix())
	sign := fmt.Sprintf("%x", sha1.Sum(append(body, []byte(ts+secret)...)))
	header.Set("Azbit-Auth-ApiKey", key)
	header.Set("Azbit-Auth-Timestamp", ts)
	header.Set("Azbit-Auth-Sign", sign)
	return header
}
