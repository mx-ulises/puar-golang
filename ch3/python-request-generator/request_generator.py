import time
import requests
import random
from threading import Thread

SAVE_RATE = 0.8
SAVE_ENDPOINT = "http://golang-app:8080/save?coins="
SPEND_ENDPOINT = "http://golang-app:8080/spend?coins="

def send_request():
    while True:
        try:
            coins = random.randint(1, 100)
            endpoint = SAVE_ENDPOINT
            if SAVE_RATE < random.random():
                endpoint = SPEND_ENDPOINT
            endpoint += str(coins)
            response = requests.get(endpoint)
            print(f"Request {endpoint}. Status Code:", response.status_code)
        except requests.exceptions.RequestException as e:
            print("Error:", e)
        time.sleep(0.1)  # 100 milliseconds is 0.1 seconds

if __name__ == "__main__":
    send_request()
