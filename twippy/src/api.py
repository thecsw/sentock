# flask_web/app.py
from flask import Flask, request, jsonify, render_template, send_file
from datetime import datetime
import matplotlib.pyplot as plt
import matplotlib.dates as pltdates
import requests
import time
from scipy import signal
import sys
import json
import databa
import psycopg2
import time
import mpld3
from mpld3 import plugins

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
databa.create_table()

#webpage asks for all stocks here
# ask === GET
# Example:
# http://localhost/api/sentiments?company=Chipotle gives all chipotle data
# http://localhost/api/sentiments givel all data
@app.route('/api/sentiments', methods=["GET"])
def get_stocks():
    company = request.args.get('company', default = '', type = str)
    before = request.args.get('before', default = int(time.time()), type = int)
    after = request.args.get('after', default = int(time.time())-12*3600, type = int)
    if (company == ''):
        return jsonify({"Error":"No company name provided!"})
    return jsonify(databa.get_sentiments(company, before, after))

@app.route('/api/graphs', methods=["GET"])
def get_graph_24h():
    hours = request.args.get('h', default = 12, type = int)
    before = int(time.time())
    after = before - hours*3600
    companies = ["McDonalds", "Fedex", "Chipotle", "Microsoft", "Disney", "Tesla", "Twitter", "Google", "Facebook", "Amazon"]
    fig, ax = plt.subplots()
    ax.set_xlabel("Date")
    ax.set_ylabel("Sentiment")
    ax.set_title("Sentimental stocks")
    for company in companies:
        values = databa.get_sentiments(company, before, after)
        x_vals = [datetime.fromtimestamp(x[0]) for x in values]
        y_vals = [x[1] for x in values]
        ax.plot_date(pltdates.date2num(x_vals), signal.medfilt(y_vals), 
        label=company, linestyle='solid', marker=None)
    ax.legend()
    handles, labels = ax.get_legend_handles_labels() # return lines and labels
    interactive_legend = plugins.InteractiveLegendPlugin(zip(handles,
                                                         ax.collections),
                                                        labels,
                                                        alpha_unsel=0.5,
                                                        alpha_over=1.5, 
                                                        start_visible=True)
    plugins.connect(fig, interactive_legend)
    return mpld3.fig_to_html(fig, figid="mygraph")

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3000)
