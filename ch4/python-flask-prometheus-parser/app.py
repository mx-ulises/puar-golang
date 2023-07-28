import os
import requests

from flask import Flask, make_response
from prometheus_client.parser import text_string_to_metric_families

app = Flask(__name__)
PUSHGATEWAY_URL = "http://pushgateway:9091/metrics"


def get_request_text(url: str) -> str:
    try:
        response = requests.get(url)
        response.raise_for_status()  # Check for any errors in the response
        content = response.text
        return content
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return None


@app.route('/')
def home():
    output = []
    for family in text_string_to_metric_families(get_request_text(PUSHGATEWAY_URL)):
        for sample in family.samples:
          output.append("Name: {0} Labels: {1} Value: {2}".format(*sample))
    content = "Metrics Families:\n" + "\n".join(output)
    response = make_response(content)
    response.headers["Content-Type"] = "text/plain"
    return response


if __name__ == '__main__':
    port = int(os.environ.get('PORT', 8080))
    app.run(debug=True, host='0.0.0.0', port=port)
