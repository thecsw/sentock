# this module handles DB actions

def init_companies(conn):
    c = conn.cursor()
    sql_query = """
    CREATE TABLE IF NOT EXISTS companies (
    rowid SERIAL PRIMARY KEY,
    company TEXT,
    timestamp INTEGER,
    sentiment REAL);"""
    c.execute(sql_query)
    conn.commit()
    c.close()

def add_sentiment(conn, company, timestamp, sentiment):
    c = conn.cursor()
    sql_query = """
    INSERT INTO companies (company, timestamp, sentiment)
    VALUES (%s, %s, %s);"""
    c.execute(sql_query, (company, timestamp, sentiment,))
    conn.commit()
    c.close()
