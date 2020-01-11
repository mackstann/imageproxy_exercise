# imageproxy_exercise

Code exercise: An http proxy that converts images to grayscale

# How to Build and Run

This project uses Python 3, so be careful not to accidentally use Python 2 if you have it installed.

    virtualenv env # or maybe virtualenv-3
    source env/bin/activate
    pip install -r requirements.txt # or maybe pip3
    IMAGEPROXY_BASE=https://maps.wikimedia.org/ FLASK_ENV=development FLASK_APP=imageproxy.py flask run

Now, find an image you like at maps.wikimedia.org, and replace the scheme/host with `http://localhost:5000`:

* https://maps.wikimedia.org/osm-intl/4/4/6@2x.png -> http://localhost:5000/osm-intl/4/4/6@2x.png
* https://maps.wikimedia.org/osm-intl/4/4/7@2x.png -> http://localhost:5000/osm-intl/4/4/7@2x.png

Non-supported file types should return appropriate HTTP errors:

* http://localhost:5000/main.css -> 501 Not Implemented

A timeout should return an HTTP 502 Bad Gateway. Set the `IMAGEPROXY_TIMEOUT` env var to `0.001` to exercise this.
