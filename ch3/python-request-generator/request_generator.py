import os
import random
import requests
import time

ENDPOINT_PATTERN = "http://{server}:{port}/v1/{path}?coins="

def send_request(server: str, port: str, save_rate: float):
    save_endpoint = ENDPOINT_PATTERN.format(server=server, port=port, path="save")
    spend_endpoint = ENDPOINT_PATTERN.format(server=server, port=port, path="spend")
    while True:
        try:
            coins = random.randint(1, 100)
            endpoint = save_endpoint
            if save_rate < random.random():
                endpoint = spend_endpoint
            endpoint += str(coins)
            response = requests.get(endpoint)
            print(f"Request {endpoint}. Status Code:", response.status_code)
        except requests.exceptions.RequestException as e:
            print("Error:", e)
        time.sleep(0.1)  # 100 milliseconds is 0.1 seconds

if __name__ == "__main__":
    server = os.getenv("TARGET_SERVER", "golang-app")
    port = os.getenv("TARGET_PORT", "8080")
    save_rate = float(os.getenv("SAVE_RATE", "0.5"))
    send_request(server, port, save_rate)
