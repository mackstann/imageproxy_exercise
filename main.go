package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
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

	typ := "image/png"
	if respType, ok := resp.Header["Content-Type"]; ok {
		typ = respType[0]
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
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
		return
	}
}

func main() {
	http.ListenAndServe(":5000", new(ImageProxyHandler))
}
