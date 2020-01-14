package main

import (
	"net/http"
)

type ImageProxyHandler struct{}

func (ImageProxyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	res.Write([]byte(path))
}

func main() {
	http.ListenAndServe(":5000", new(ImageProxyHandler))
}
