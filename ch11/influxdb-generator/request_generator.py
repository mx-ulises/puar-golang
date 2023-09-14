import os
import random
import requests
import time

ENDPOINT_PATTERN = "http://{server}:{port}/write"
DATA_PATTERN = 'example_metric,tagname=value value={value}'

def send_request(server: str, port: str, positive_rate: float):
    value = 0
    endpoint = ENDPOINT_PATTERN.format(server=server, port=port)
    while True:
        try:
            if random.random() < positive_rate:
                value += 1
            else:
                value -= 1
            data = DATA_PATTERN.format(value=value)
            response = requests.post(
                endpoint,
                data=data.encode('utf-8'),
                headers={'Content-Type': 'application/octet-stream'}
            )
            print(f"Request {endpoint}. Status Code:", response.status_code)
        except requests.exceptions.RequestException as e:
            print("Error:", e)
        time.sleep(0.1)  # 100 milliseconds is 0.1 seconds

if __name__ == "__main__":
    server = os.getenv("TARGET_SERVER", "golang-app")
    port = os.getenv("TARGET_PORT", "8080")
    save_rate = float(os.getenv("POSITIVE_RATE", "0.5"))
    send_request(server, port, save_rate)
