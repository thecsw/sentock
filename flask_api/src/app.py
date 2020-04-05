# flask_web/app.py
from flask import Flask, request, jsonify, render_template
import sys
import json

app = Flask(__name__)

#this is where the webpage is going to pull from. It's going
#to always pull from here for data. I can't think of anything
#else that the webpage is going to need for mvp

#webpage asks for all stocks here
@app.route('/api/stocks', methods=["POST"])
def get_stocks():
    raw_data = request.data
    stock_name = json.loads(raw_data)['stock_name']
    print(stock_name, file=sys.stderr)
    return "hello there"

#static page for docs of this
@app.route('/api/docs', methods=["GET"])
def serve():
    return render_template('app.html')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3000)