import requests
from bs4 import BeautifulSoup
import re
import smtplib
import time
from datetime import datetime
import os


MY_EMAIL = os.getenv("MY_EMAIL")
PASSWORD = os.getenv("PASSWORD")
RECIPIENTS = os.getenv("RECIPIENTS")
URL = os.getenv("URL")


dict = {}

data = requests.get(URL)

dict1 = data.json()["sheet1"]

for items in dict1:
    dict[items['bookName']] = {
        items['bookSearch']:items["pricePoint"]
    }

start_time = time.time()
for books in dict:
    for items in dict[books]:
        url = items
        page = requests.get(URL)
        data = BeautifulSoup(page.content, "html.parser")
        all = data.find_all("tr", class_=['results-table-first-LogoRow has-data', "results-table-LogoRow has-data"])
        for words in all:
            if "fair" in words.text.lower():
                break
            prices = words.find_all("span", class_="results-price")
            for nums in prices:
                full_prices = nums.find("a").text
                new = re.sub('\D', '', full_prices)[:-2]
                print(f"{books} : {new}")
                if int(new) < int(dict[books][items]):
                    id = list(dict).index(books) + 2
                    print(id)
                    data = {
                        "sheet1": {
                            "pricePoint":int(new)
                        }
                    }
                    response = requests.put(url=f"{URL}/{id}", json=data)
                    print(response.status_code)
                    with smtplib.SMTP("smtp.gmail.com") as connection:
                        connection.starttls()
                        connection.login(user=MY_EMAIL, password=PASSWORD)
                        connection.sendmail(from_addr=MY_EMAIL,
                                    to_addrs=RECIPIENTS,
                                    msg=f"Subject: Book Alert for {books} \n\n Ander, {books} has dropped to ${new}")
                        break
            
print(f"Cycle completed at {datetime.now()}")
print(f"{(time.time() - start_time)} seconds to complete")
