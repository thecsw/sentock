import time
import psycopg2

# Sleep for 3 seconds before DB starts accepting
# connections
time.sleep(3)
# Initialize the database
conn = psycopg2.connect(dbname="sentock",
                        user="sandy",
                        password="pass",
                        host="postgres",
                        port="5432")

def create_table():
    c = conn.cursor()
    c.execute("""
    CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    company TEXT,
    tweet_id TEXT UNIQUE,
    original TEXT,
    timestamp INT,
    sentiment REAL);""")
    conn.commit()
    c.close()

def add_sentiment(company, tweet_id, original, timestamp, sentiment):
    c = conn.cursor()
    c.execute("""
    INSERT INTO companies (company, tweet_id, original, timestamp, sentiment)
    VALUES (%s, %s, %s, %s, %s)
    """, (company, tweet_id, original, timestamp, sentiment, ))
    conn.commit()
    c.close()

def get_sentiments(company, before, after):
    c = conn.cursor()
    c.execute("""
    SELECT timestamp, sentiment FROM companies WHERE company = %s AND timestamp < %s AND timestamp > %s ORDER BY timestamp DESC;
    """, (company, int(before), int(after),))
    result = c.fetchall()
    c.close()
    if (len(result) < 1):
        return []
    #moving average (of time period of size [window] seconds):
    window = 60*60 #window size of an hour
    latest = (result[0][0]//60)*60 #anchor on the minute mark
    end_result = []
    start = 0
    while latest > result[len(result)-1][0]:
        print(before, start, flush=True)
        (window_ave, start) = get_window_ave(result, window, latest, start)
        if window_ave is not None:
            end_result.append((latest,window_ave))
        latest-=60
    return end_result

def close():
    conn.close()

def get_window_ave(vals, window, latest, start):
    (values, startNew) = get_vals_of_interval(vals, latest, latest-window, start)
    if (len(values)<=5):
        return (None, startNew)
    sumVal = 0 
    for val in values:
        sumVal += val[1]
    return (sumVal / len(values), startNew)

def get_vals_of_interval(vals, before, after, start):
    ans = []
    startNew = start
    for val in vals[start:]:
        if (val[0] > before - 60):
            startNew += 1
        if (val[0] < before and val[0] > after):
            ans.append(val)
        if (val[0] < after):
            break
    return (ans, startNew)