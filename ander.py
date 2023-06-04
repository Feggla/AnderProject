import requests
from bs4 import BeautifulSoup
import re
import smtplib
import time
from datetime import datetime
import asyncio
import aiohttp
import html5lib


MY_EMAIL = "feggylad@gmail.com"
PASSWORD = "sltyefjbswqlvlox"
RECIPIENTS = ["michaelfeggans@gmail.com"]

dict = {
    "WoT 1" : {
        "https://www.bookfinder.com/search/?keywords=0312850093&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 30
    },
    "WoT 2" : {
        "https://www.bookfinder.com/search/?keywords=0312851405&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 30
    },
    "WoT 3" : {
        "https://www.bookfinder.com/search/?keywords=0312852487&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 30
    },
    "WoT 4" : {
        "https://www.bookfinder.com/search/?keywords=0312854315&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 30
    },
    "WoT 5" : {
        "https://www.bookfinder.com/search/?keywords=0312854277&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 30
    },
    "WoT 6" : {
        "https://www.bookfinder.com/search/?keywords=0312854285&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 25
    },
    "WoT 7" : {
        "https://www.bookfinder.com/search/?keywords=0312857675&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 25
    },
    "WoT 8" : {
        "https://www.bookfinder.com/search/?keywords=0312857691&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 25
    },
    # "WoT 9" : {
    #     "https://www.bookfinder.com/search/?keywords=0312864256&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 25
    # },
    "WoT 10" : {
        "https://www.bookfinder.com/search/?keywords=0312864590&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 25
    },
    "WoT 11" : {
        "https://www.bookfinder.com/search/?keywords=0312873077&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 25
    },
    "WoT 12" : {
        "https://www.bookfinder.com/search/?keywords=0765302306&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 20
    },
    "WoT 13" : {
        "https://www.bookfinder.com/search/?keywords=0765325942&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 20
    },
    "WoT 14" : {
        "https://www.bookfinder.com/search/?keywords=0765325950&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=" : 20
    },
    "Mistress of Empire" : {
        "https://www.bookfinder.com/search/?keywords=0246133554&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=&fbclid=IwAR1sMWaf_iTY-1G2ZMD4yTrdydx-Ecy-E9g7VZSpHYNfeqtDaP-6YMEotHY" : 39
    },
    "Magician" : {
        "https://www.bookfinder.com/search/?keywords=0246122056&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=&fbclid=IwAR3_yA_uDDY_LdFN24jlDhMc0Dk4wP3UVrMiGvLo9xLaXKHxnTeInfQlbdg" : 0
    },
    "Silverthorn" : {
        "https://www.bookfinder.com/search/?keywords=0246125411&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=&fbclid=IwAR2kbOToxIetdTR2RzKaStpQ3Pf6vUhH0PuHCWxCzQsHkQ7UyZ0jC0YccEw" : 150
    },
    "Sethanon" : {
        "https://www.bookfinder.com/search/?keywords=0246128283&currency=AUD&destination=au&mode=advanced&il=en&classic=off&lang=en&st=sh&ac=qr&submit=&fbclid=IwAR1U_oYHwAxN-NbzUenPyh8yBbdZNDzJM1AQl8nfjLcJchQBEJX50Mza6yA" : 150
    }
}

async def main():

    async with aiohttp.ClientSession() as session:
        start_time = time.time()
        for books in dict:
                for items in dict[books]:
                    url = items
                    async with session.get(url) as resp:
                        page = await resp.text()
                        data = BeautifulSoup(page, "html.parser")
                        prices = data.find_all("span", class_="results-price")

                        for nums in prices:
                                full_prices = nums.find("a").text
                                new = re.sub('\D', '', full_prices)[:-2]
                                # print(f"{books} : {new}")
                                print(f"{new}:{dict[books][items]}")
                                if int(new) <= int(dict[books][items]):
                                    with smtplib.SMTP("smtp.gmail.com") as connection:
                                        connection.starttls()
                                        connection.login(user=MY_EMAIL, password=PASSWORD)
                                        connection.sendmail(from_addr=MY_EMAIL,
                                                to_addrs=RECIPIENTS,
                                                msg=f"Subject: Book Alert for {books} \n\n Ander, {books} has dropped to ${new}")
        print(f"Cycle completed at {datetime.now()}")
        print(f"{(time.time() - start_time)} seconds to complete")

asyncio.run(main())


# while True:
#     start_time = time.time()
#     for books in dict:
#         for items in dict[books]:
#             url = items
#             page = requests.get(url)
#             data = BeautifulSoup(page.content, "html.parser")

#             prices = data.find_all("span", class_="results-price")

#             for nums in prices:
#                 full_prices = nums.find("a").text
#                 new = re.sub('\D', '', full_prices)[:-2]
#                 if int(new) <= int(dict[books][items]):
#                     with smtplib.SMTP("smtp.gmail.com") as connection:
#                         connection.starttls()
#                         connection.login(user=MY_EMAIL, password=PASSWORD)
#                         connection.sendmail(from_addr=MY_EMAIL,
#                                     to_addrs=RECIPIENTS,
#                                     msg=f"Subject: Book Alert for {books} \n\n Ander, {books} has dropped to ${new}")
#                 break
#     print(f"Cycle completed at {datetime.now()}")
#     print(f"{(time.time() - start_time)} seconds to complete")
#     time.sleep(60)
