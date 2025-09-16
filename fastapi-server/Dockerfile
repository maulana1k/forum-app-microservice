# -----------------------------
# Build stage
# -----------------------------
FROM python:3.12-slim AS builder

WORKDIR /fastapi-server

# Install build dependencies (for packages like numpy, psycopg2)
RUN apt-get update && apt-get install -y build-essential gcc libpq-dev --no-install-recommends

# Install Python dependencies into a temporary directory
COPY requirements.txt .
RUN pip install --prefix=/install --no-cache-dir -r requirements.txt

# Copy source
COPY ./src ./src
COPY .env .env

# -----------------------------
# Run stage (minimal)
# -----------------------------
FROM python:3.12-slim

WORKDIR /fastapi-server

# Copy installed packages
COPY --from=builder /install /usr/local

# Copy source
COPY ./src ./src
COPY .env .env

# Expose ports
EXPOSE 8000 50051

# Run FastAPI
CMD ["python", "-m", "src.main"]
