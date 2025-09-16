FROM apache/airflow:2.7.1-python3.12

# Set environment variables
ENV AIRFLOW_HOME=/opt/airflow
WORKDIR ${AIRFLOW_HOME}

# Install extra dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy DAGs and scripts
COPY dags/ ${AIRFLOW_HOME}/dags/
COPY scripts/ ${AIRFLOW_HOME}/scripts/

# Set entrypoint
# For production/dev, we will use docker-compose to run scheduler/webserver separately
