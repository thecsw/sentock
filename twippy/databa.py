import sqlite3

conn = sqlite3.connect("sentocks.db", check_same_thread=False)

def create_table():
    c = conn.cursor()
    c.execute("""
    CREATE TABLE IF NOT EXISTS companies (
    id INTEGER PRIMARY KEY,
    company TEXT,
    tweet_id TEXT UNIQUE,
    original TEXT,
    timestamp INTEGER,
    sentiment REAL);""")
    conn.commit()
    c.close()

def add_sentiment(company, tweet_id, original, timestamp, sentiment):
    c = conn.cursor()
    c.execute("""
    INSERT INTO companies (company, tweet_id, original, timestamp, sentiment)
    VALUES (?, ?, ?, ?, ?)
    """, (company, tweet_id, original, timestamp, sentiment, ))
    conn.commit()
    c.close()

def close():
    conn.close()