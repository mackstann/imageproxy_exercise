package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ImageProxyHandler struct {
	Timeout time.Duration
	Base    string
}

func (handler *ImageProxyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	client := &http.Client{Timeout: handler.Timeout}

	resp, err := client.Get(handler.Base + "/" + path)
	if err != nil {
		if tErr, ok := err.(net.Error); ok && tErr.Timeout() {
			log.Printf("Timeout while proxying path %q", path)
			res.WriteHeader(502) // 503 might be more correct?
			return
		}

		// 500 is never ideal, but without more specific error handling, we don't know exactly what kind of
		// problem happened.
		log.Printf("Error %q while proxying %q", err.Error(), path)
		res.WriteHeader(500)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Printf("HTTP %d while proxying %q", resp.StatusCode, path)
		res.WriteHeader(resp.StatusCode) // proxy HTTP errors
		return
	}

	typ := "image/png"
	if respType, ok := resp.Header["Content-Type"]; ok {
		typ = respType[0]
	}

	switch typ {
	case "image/png", "image/jpeg":
	default:
		log.Printf("Content-Type %q not supported", typ)
		res.WriteHeader(501) // 502 might be more correct?
		return
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("Failed to decode image: %s", err.Error())
		res.WriteHeader(500)
		return
	}

	// convert to grayscale
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	encodeFunc := png.Encode
	if typ == "image/jpeg" {
		encodeFunc = func(w io.Writer, im image.Image) error {
			return jpeg.Encode(w, im, nil)
		}
	}

	res.Header().Set("Content-Type", typ)

	if err := encodeFunc(res, grayImg); err != nil {
		log.Printf("Failed to encode image: %s", err.Error())
		res.WriteHeader(500)
		return
	}
}

func main() {
	handler := new(ImageProxyHandler)

	timeoutMS := 5000
	if timeoutStr := os.Getenv("IMAGEPROXY_TIMEOUT"); timeoutStr != "" {
		var err error
		timeoutMS, err = strconv.Atoi(timeoutStr)
		if err != nil {
			log.Fatalf("Please set env var IMAGEPROXY_TIMEOUT to a valid number of milliseconds")
		}
	}
	handler.Timeout = time.Duration(timeoutMS) * time.Millisecond

	base := os.Getenv("IMAGEPROXY_BASE")
	if base == "" {
		log.Fatalf("Please set env var IMAGEPROXY_BASE to a valid URL")
	}
	base = strings.TrimRight(base, "/")
	handler.Base = base

	http.ListenAndServe(":5000", handler)
}
