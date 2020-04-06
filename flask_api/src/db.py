def get_company(conn, company):
    c = conn.cursor()
    sql_query = """
    SELECT * FROM companies WHERE company = %s;"""
    c.execute(sql_query, (company,))
    results = c.fetchall()
    c.close()
    return results

def get_all(conn):
    c = conn.cursor()
    sql_query = """
    SELECT * FROM companies;"""
    c.execute(sql_query)
    results = c.fetchall()
    c.close()
    return results
