# 🗨️ **Forum Microservices App**

> **AI-powered discussion platform** built with **Next.js, GoFiber, FastAPI, gRPC, RabbitMQ, PostgreSQL, Redis**, and **Docker Compose**.

A production-style forum where users can create posts, comment, and receive real-time sentiment analysis using a Transformer model.
This project demonstrates **full-stack microservice engineering**—from API design and AI inference to load testing and monitoring.

---

## ✨ Features

* **🔒 Authentication & Authorization** – GoFiber service with PostgreSQL & JWT.
* **📝 Posts Service** – Create, read, and manage forum posts with Redis caching.
* **🤖 AI Inference** – FastAPI service running a HuggingFace Transformer (ONNX export) for sentiment & content moderation.
* **📬 Async Messaging** – RabbitMQ for decoupled communication between services.
* **⚡ gRPC** – High-performance RPC between Go and Python services.
* **📈 Monitoring & Stress Testing** – Grafana + k6 for metrics and load testing.
* **🐳 Fully Containerized** – One-command spin-up with Docker Compose.

---
## 🗺️ Architecture Overview

```text
         ┌──────────────────┐
         │      Client      │
         │ (Next.js Web)    │
         └───────▲──────────┘
                 │ REST
                 │
        ┌────────┴─────────┐
        │   Go Server      │
        │ (Fiber + gRPC)   │
        └─▲──────┬─────────┘
  Feeds API │    │ Create Post
   gRPC     │    │ Publish
            │    │
      ┌─────┘    ▼
┌──────────────┐   ┌────────────────────────┐
│  FastAPI ML  │   │   RabbitMQ Broker      │
│ (Inference)  │◄──┤  (Async Messaging)     │
└──────────────┘   └────────────────────────┘
        ▲
        │ gRPC Response
        │
 ┌───────┴────────┐
 │ Postgres /     │
 │ Redis Cache    │
 └────────────────┘
```

## 🚀 Quick Start

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/maulana1k/forum-app-microservice.git
cd forum-app-microservice
```

### 2️⃣ Environment Setup

Create a `.env` file in the root (sample below):

```env
POSTGRES_USER=dev
POSTGRES_PASSWORD=dev
POSTGRES_DB=forumdb
POSTGRES_PORT=5432
POSTGRES_HOST=postgres

REDIS_HOST=redis
REDIS_PORT=6379

RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672

DOCKER_ENV=true
```

> 💡 **Tip:** Adjust values to match your local/dev environment.

### 3️⃣ Build & Run with Docker Compose

```bash
docker-compose up --build
```

This will start:

* **Next.js Web Client** (frontend)
* **GoFiber Auth Service**
* **GoFiber Posts Service**
* **FastAPI AI Inference Service**
* **PostgreSQL**, **Redis**, and **RabbitMQ**

Visit:

* **Frontend:** [http://localhost:3000](http://localhost:3000)
* **Golang Service:** [http://localhost:8080](http://localhost:8080)
* **Fastapi Service:** [http://localhost:8000](http://localhost:8000)
* **Grafana:** [http://localhost:3001](http://localhost:3001)
* **RabbitMQ Management:** [http://localhost:15672](http://localhost:15672) (guest/guest)

---

## 🧪 Demo & Testing

### 🔹 Send Bulk Test Messages

```bash
python scripts/send_bulk_messages.py
```

Publishes sample posts to RabbitMQ for inference & moderation.

### 🔹 Stress Test with k6

```bash
k6 run k6/stress_test.js
```

Simulates high traffic to measure throughput & latency.

---

## 🛠️ Development Notes

* **Hot Reload:** Frontend and backend services support live reload during development.
* **Database Migrations:** Managed via `go-migrate` (see `/migrations`).
* **Model Updates:** Retrain and export ONNX models using `fastapi-server/notebooks/sentiment/model.ipnyb`

---

## 🧰 Tech Stack

* **Frontend:** [Next.js](https://nextjs.org/)
* **Backend:** [GoFiber](https://gofiber.io/), [FastAPI](https://fastapi.tiangolo.com/)
* **Database:** PostgreSQL + Redis
* **Messaging:** RabbitMQ
* **Inter-service:** gRPC
* **Containerization:** Docker & Docker Compose
* **Monitoring:** Grafana + k6
* **Machine Learning:** HuggingFace Transformers (PyTorch → ONNX)

---

## 🌱 Contributing

Pull requests are welcome!
Please open an issue to discuss major changes or improvements.

---

