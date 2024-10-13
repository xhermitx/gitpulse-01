# GitPulse Tracker

**GitPulse Tracker** is a streamlined recruiting tool designed to optimize talent acquisition processes for tech companies. By integrating resume parsing, GitHub data analysis, and candidate profiling, the system provides recruiters with insights into candidates' technical capabilities.

## Architecture

The architecture of GitPulse Tracker revolves around several core services working together:

![GitPulse Architecture](./static/architecture.png "GitPulse Architecture")

- **Drive**: Handles resume file inputs for processing.
- **Resume-Parser**: Extracts candidate information from resumes and fetches their GitHub profiles.
- **Redis**: Used as a cache layer between services.
- **RabbitMQ**: Manages asynchronous tasks across services.
- **Profiler**: Analyzes candidate GitHub data to compute a technical score.
- **MySQL**: Stores user and job-related data.
- **GitPulse (UI)**: Provides the interface for recruiters to interact with the system.

## Endpoints

The system exposes a set of RESTful APIs to interact with GitPulse Tracker services. Swagger UI provides detailed documentation and testing capabilities for the available endpoints:

![API Specifications](./static/endpoints.png "API Specifications")

### Key API Categories:
- **Auth**: User registration, login, and account management.
- **Job**: Job creation, status updates, and job result fetching.
- **Example**: Example API to demonstrate basic functionality.

## Project Setup

### 1. Prerequisites

GitPulse Tracker leverages several external dependencies:
- **RabbitMQ**
- **Redis**
- **MySQL**

Ensure Docker and Docker Compose are installed on your system to manage these services as containers.

### 2. Docker Setup

The Docker environment manages RabbitMQ, Redis, and MySQL services. The `docker-compose.yml` file includes the necessary configurations. You need to add your environment-specific settings in a `.env` file. The variables are referenced in the `docker-compose.yml`.

- To start the containers:
  ```bash
  docker-compose up -d
  ```

- To stop the containers:
  ```bash
  docker-compose down
  ```

Ensure that the containers are running properly before proceeding.

### 3. Service Environment Configuration

The project is divided into three main services:
- **Backend**: Handles API requests for user authentication, job management, and other core functionalities.
- **Resume-Parser**: Extracts GitHub usernames from resumes.
- **Profiler**: Processes GitHub data and assigns a candidate score based on predefined metrics.

Each service has a `Makefile` and a `config` folder containing an `env.go` file. You should set up the `.env` file for each service by following the guidelines in the `env.go` file.

Once the environment variables are properly configured, proceed to the next step.

### 4. Running Migrations

Migrations are used to set up the MySQL database schema for the backend. Ensure the MySQL instance is running in the Docker container.

From the `./backend` directory, you can run the following command to apply the database migrations:
```bash
make migrate-up
```

This will create the necessary tables and schema in your database.

### 5. Running the Services

Once all environment variables are configured, and migrations have been applied, you can start the services. Use the following commands to build and run each service:

- **Backend**:
  ```bash
  cd ./backend
  make build
  make run
  ```

- **Resume-Parser**:
  ```bash
  cd ./resume-parser
  make build
  make run
  ```

- **Profiler**:
  ```bash
  cd ./profiler
  make build
  make run
  ```

Each service should now be running successfully. Ensure all services are able to connect to the database, Redis, and RabbitMQ.

### 6. API Documentation

To explore and interact with the APIs, navigate to the Swagger UI provided by the backend:

[Swagger UI](http://localhost:8000/swagger/index.html#/)

Swagger provides a user-friendly interface to test the available endpoints for authentication, job management, and more.