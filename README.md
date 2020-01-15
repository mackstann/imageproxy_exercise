# imageproxy_exercise

Code exercise: An http proxy that converts images to grayscale

# How to Build and Run

This project uses Go and its standard library only; no third-party packages.

    go run main.go

Now, find an image you like at maps.wikimedia.org, and replace the scheme/host with `http://localhost:5000`:

* https://maps.wikimedia.org/osm-intl/4/4/6@2x.png -> http://localhost:5000/osm-intl/4/4/6@2x.png
* https://maps.wikimedia.org/osm-intl/4/4/7@2x.png -> http://localhost:5000/osm-intl/4/4/7@2x.png

Non-supported file types should return appropriate HTTP errors:

* http://localhost:5000/main.css -> 501 Not Implemented

A timeout should return an HTTP 502 Bad Gateway. Set the `IMAGEPROXY_TIMEOUT` env var to `1` to exercise this. (The unit
is milliseconds.)
