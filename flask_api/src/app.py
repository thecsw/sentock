# flask_web/app.py
from flask import Flask, request, jsonify
import sys
import json

app = Flask(__name__)

#this is where the webpage is going to pull from. It's going 
#to always pull from here for data. I can't think of anything
#else that the webpage is going to need for mvp

#webpage sends {"stock_name":"<stockname>"}
@app.route('/', methods=["POST"])
def hello_world():
    raw_data = request.data
    stock_name = json.loads(raw_data)['stock_name']
    print(stock_name, file=sys.stderr)

    return "hello there"

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=3000)