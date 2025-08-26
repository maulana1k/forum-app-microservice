FROM python:3.12-slim

WORKDIR /fastapi-server

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY ./src ./src
COPY .env .env

EXPOSE 8000
CMD ["python", "-m", "src.main"]
