package main

import (
	"image"
	_ "image/jpeg"
	"image/png"
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

	if err := png.Encode(res, grayImg); err != nil {
		return
	}
}

func main() {
	http.ListenAndServe(":5000", new(ImageProxyHandler))
}
