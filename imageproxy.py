import os, sys, logging, socket, io
from urllib.request import urlopen
from urllib.error import HTTPError, URLError

from PIL import Image

from flask import Flask, request, abort, url_for
from flask.json import jsonify

app = Flask(__name__)

TIMEOUT = 5

try:
    BASE = os.environ['IMAGEPROXY_BASE'].rstrip('/')
except KeyError:
    logging.error('Please set env var IMAGEPROXY_BASE to a valid URL')
    os._exit(1)

@app.route('/', defaults={'path': ''}) # flask needs a special case for root path
@app.route('/<path:path>')
def proxy(path):
    content_type = None
    body = None

    try:
        with urlopen(BASE + '/' + path, timeout=TIMEOUT) as response:
            content_type = response.headers.get('Content-Type', '<none>')
            body = response.read()
    except HTTPError as err:
        return abort(err.code)
    except ValueError as e:
        logging.error("Can't proxy invalid path %r", path)
        return abort(400)
    except URLError as err:
        if isinstance(err.reason, socket.timeout):
            logging.error('Timeout while proxying path %r', path)
            return abort(502) # 503 might be more correct?
        else:
            raise

    # In real life it might be advantageous to convert the image to another type, and although it might be a bit
    # misleading at first glance (compared to the filename in the URL), as long as the client accepts that mime type
    # (i.e. proper content negotiation), there's nothing wrong with it per se.
    #
    # But per the requirements (and it is of course simpler anyway), this maintains the same image type.

    type_map = {
        'image/png': 'png',
        'image/jpeg': 'jpeg',
        'image/gif': 'gif',
    }

    if content_type not in type_map:
        # With more time I'd work on backup strategies such as sniffing the image data to detect the type.
        logging.error("Content-Type %s not supported", content_type)
        return abort(501) # 502 might be more correct?

    return (
        build_response(type_map[content_type], body),
        200,
        {'Content-Type': content_type},
    )

def build_response(format, imgdata):
    img = Image.open(io.BytesIO(imgdata)).convert('L') # L means grayscale

    with io.BytesIO() as output:
        img.save(output, format=format)
        return output.getvalue()
