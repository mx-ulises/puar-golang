FROM python:3.9-slim-buster

#WORKDIR /app

#COPY requirements.txt .

RUN pip install requests

COPY . .

CMD ["python", "request_generator.py"]