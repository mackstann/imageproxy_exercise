# imageproxy_exercise

Code exercise: An http proxy that converts images to grayscale

# How to Build and Run

This project uses Python 3, so be careful not to accidentally use Python 2 if you have it installed.

    virtualenv env # or maybe virtualenv-3
    source env/bin/activate
    pip install -r requirements.txt # or maybe pip3
    FLASK_ENV=development FLASK_APP=imageproxy.py flask run
