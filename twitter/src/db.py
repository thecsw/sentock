# this module handles DB actions

def init_company(conn, company):
    # Quick escaping the name
    company = company.replace("'", "")
    c = conn.cursor()
    # Doing a hack with concatenating because
    # psycopg2 doesn't allow variable table names
    sql_query = """
    CREATE TABLE IF NOT EXISTS TABLE_NAME (
    rowid SERIAL PRIMARY KEY,
    timestamp INTEGER,
    sentiment REAL);""".replace("TABLE_NAME", company)
    c.execute(sql_query)
    conn.commit()
    c.close()

def add_sentiment(conn, company, timestamp, sentiment):
    c = conn.cursor()
    sql_query = """
    INSERT INTO %s (timestamp, sentiment)
    VALUES (%s, %s);"""
    c.execute(sql_query, (company, timestamp, sentiment,))
    conn.commit()
    c.close()
