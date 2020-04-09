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
    return result

def close():
    conn.close()
