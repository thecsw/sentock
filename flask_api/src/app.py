# flask_web/app.py
from flask import Flask, request, jsonify, render_template
import sys
import json
import db
import psycopg2
import time

app = Flask(__name__)

# Sleep for 3 seconds before DB starts accepting
# connections
time.sleep(3)
# Initialize the database
conn = psycopg2.connect(dbname="sentock",
                        user="sandy",
                        password="pass",
                        host="postgres",
                        port="5432")

#webpage asks for all stocks here
# ask === GET
# Example:
# http://localhost/api/sentiments?company=Chipotle gives all chipotle data
# http://localhost/api/sentiments givel all data
@app.route('/api/sentiments', methods=["GET"])
def get_stocks():
    company = request.args.get('company', default = '', type = str)
    if (company == ''):
        return jsonify(db.get_all(conn))
    return jsonify(db.get_company(conn, company))

#static page for docs of this
@app.route('/api/docs', methods=["GET"])
def serve():
    return render_template('app.html')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3000)
