package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

type ImageProxyHandler struct{}

func (ImageProxyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get("https://maps.wikimedia.org" + path)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	res.Write(body)

}

func main() {
	http.ListenAndServe(":5000", new(ImageProxyHandler))
}
