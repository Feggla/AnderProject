import gspread
from oauth2client.service_account import ServiceAccountCredentials
import requests
from bs4 import BeautifulSoup
import re
from smtplib import SMTP_SSL as SMTP
import time
from datetime import datetime
import os
import pygsheets
from dotenv import load_dotenv

load_dotenv()

MY_EMAIL = os.getenv("MY_EMAIL")
PASSWORD = os.getenv("PASSWORD")
RECIPIENT1 = os.getenv("RECIPIENT1")
RECIPIENT2 = os.getenv("RECIPIENT2")


gc = pygsheets.authorize(service_file="./ander-project-388823-d897fb774f87.json")

ws = gc.open('Ander-project')[0]
dictlist = ws.get_all_records()


# for items in dictlist:
#     print(type(items))

# for items in sheet:
# data = requests.get(URL)
# print(data.json())
# dict1 = data.json()

dict = {}

for items in dictlist:
    dict[items['Book Name']] = {
        items['Book Search']:items["Price Point"]
    }

start_time = time.time()
for books in dict:
    for items in dict[books]:
        link = items
        page = requests.get(link)
        data = BeautifulSoup(page.content, "html.parser")
        all = data.find_all("tr", class_=['results-table-first-LogoRow has-data', "results-table-LogoRow has-data"])
        for words in all:
            if "fair" in words.text.lower():
                print("Fair found")
                break
            prices = words.find_all("span", class_="results-price")
            for nums in prices:
                full_prices = nums.find("a").text
                new = re.sub('\D', '', full_prices)[:-2]
                print(f"{books} : {new}")
                if int(new) < int(dict[books][items]):
                    id = list(dict).index(books) + 2
                    print(id)
                    ws.update_value(f"C{id}", new)
                    with SMTP("smtp.gmail.com") as connection:
                        connection.set_debuglevel(False)
                        connection.login(user=MY_EMAIL, password=PASSWORD)
                        connection.sendmail(from_addr=MY_EMAIL,
                                    to_addrs=[RECIPIENT1,RECIPIENT2],
                                    msg=f"Subject: Book Alert for {books} \n\n Ander, {books} has dropped to ${new}")
                    break
            
print(f"Cycle completed at {datetime.now()}")
print(f"{(time.time() - start_time)} seconds to complete")
