# Python & Virtual Environment
PYTHON=python
PIP=$(PYTHON) -m pip
VENV?=venv

# App settings
APP=src.main:app
HOST=0.0.0.0
PORT=8000

# Default target
.DEFAULT_GOAL := help

## --- Dev Workflow ---

help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

venv: ## Create virtual environment
	$(PYTHON) -m venv $(VENV)

install: venv ## Install dependencies
	$(VENV)/Scripts/$(PIP) install --upgrade pip setuptools wheel
	$(VENV)/Scripts/$(PIP) install -r requirements.txt

mlflow: ## Run MLFlow server
	@echo "Starting MLflow server..."
	# Start MLflow in the background and redirect logs
	start /B $(VENV)/Scripts/python -m mlflow server ^
		--host 127.0.0.1 ^
		--port 5000 ^
		--backend-store-uri sqlite:///mlflow.db ^
		--default-artifact-root ./mlruns
	# Wait a few seconds for MLflow to be ready
	@echo "Waiting 5 seconds for MLflow to start..."
	@sleep 5

dev: ## Run FastAPI with auto-reload
	$(VENV)/Scripts/uvicorn $(APP) --reload --host $(HOST) --port $(PORT)

stop-dev:  ## Stop FastAPI dev server and MLflow
	@echo "Stopping MLflow and Uvicorn processes..."
	# Kill MLflow server (port 5000)
	for /f "tokens=5" %%a in ('netstat -a -n -o ^| findstr :5000') do taskkill /PID %%a /F
	# Kill Uvicorn dev server (port 8000)
	for /f "tokens=5" %%a in ('netstat -a -n -o ^| findstr :8000') do taskkill /PID %%a /F
	@echo "All dev servers stopped."


run: mlflow ## Run FastAPI in production mode
	$(VENV)/Scripts/uvicorn $(APP) --host $(HOST) --port $(PORT)

freeze: ## Freeze current dependencies to requirements.txt
	$(VENV)/Scripts/$(PIP) freeze > requirements.txt

## --- Code Quality ---

lint: ## Run linter (flake8 & black)
	$(VENV)/Scripts/$(PIP) install flake8 black isort
	$(VENV)/Scripts/flake8 .
	$(VENV)/Scripts/black --check .
	$(VENV)/Scripts/isort --check-only .

format: ## Auto-format code (black & isort)
	$(VENV)/Scripts/black .
	$(VENV)/Scripts/isort .

## --- Testing ---

test: ## Run tests with pytest
	$(VENV)/Scripts/$(PIP) install pytest
	$(VENV)/Scripts/pytest -v --maxfail=1 --disable-warnings

## --- Maintenance ---

clean: ## Remove caches, pyc files, and build artifacts
	rm -rf __pycache__ .pytest_cache *.pyc *.pyo *.pyd *.egg-info .mypy_cache

rebuild: clean install ## Clean and reinstall everything

