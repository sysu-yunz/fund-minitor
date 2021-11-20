import json
from pymongo import MongoClient
from srv import db_srv

# get data from api http://fund.eastmoney.com/js/fundcode_search.js
# get file fundcode_search.js
# fundcode_search.js to fundcode_search --- remove the first line and last line, replace [ and ] to { and }
# use this file to write to mongo

with open("fundcode_search", "r") as f, open("fundcode_search.json", "w") as n:
    count = 0
    for l in f:
        if count % 7 == 0 or count % 7 == 6:
            n.write(l)
        elif count % 7 == 1:
            n.write(f'"fundCode": ' + l)
        elif count % 7 == 2:
            n.write(f'"fundShort": ' + l)
        elif count % 7 == 3:
            n.write(f'"fundName": ' + l)
        elif count % 7 == 4:
            n.write(f'"fundType": ' + l)
        elif count % 7 == 5:
            n.write(f'"fundPinYin": ' + l)
        count = count + 1

with open("fundcode_search", "r") as f, open("fundcode_search.csv", "w") as n:
    count = 0
    for l in f:
        if count % 7 == 6:
            n.write(f'\n')
        elif count % 7 == 1:
            n.write(l)
        elif count % 7 == 2:
            n.write(l)
        elif count % 7 == 3:
            n.write(l)
        elif count % 7 == 4:
            n.write(l)
        elif count % 7 == 5:
            n.write(l)
        count = count + 1


#db_srv = "mongodb://127.0.0.1:27017"

myclient = MongoClient(db_srv, connect=False)
db = myclient.fund
mycol = db["basic"]

with open('fff_new.json') as f:
    file_data = json.load(f)

# if pymongo >= 3.0 use insert_many() for inserting many documents
mycol.insert_many(file_data)

myclient.close()
