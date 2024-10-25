# Event Booking System

This project is an event booking system with multiple services including a main application, background tasks, and a central authentication service.

## Prerequisites

- Docker and Docker Compose
- Go 1.23 or later
- Make

## Project Structure

The project consists of three main services:
1. Main Application (app)
2. Background Tasks
3. Central Authentication

## Setup Instructions

1. Clone the repository:
   ```
   git clone git@github.com:danghovu/booking.git
   cd booking-event
   ```

2. Build and start the services using Docker Compose:
   ```
   make docker/up
   ```
   This command will build the Docker images and start all services defined in the `docker-compose.yaml` file.

3. To stop and remove the containers, networks, and volumes:
   ```
   make docker/down
   ```

## Running Services Individually

If you prefer to run services individually without Docker, you can use the following commands:

1. Run the main server:
   ```
   make run-servers
   ```

2. Run the central authentication service:
   ```
   make run-central-auth
   ```

3. Run the background worker service:
    ```
    make run-worker
    ```

## Service Details

### Main Application (app)
- Dockerfile: `Dockerfile`
- Handles main business logic
- Depends on PostgreSQL and Redis

### Background Tasks
- Dockerfile: `Dockerfile.background-tasks`
- Runs background processing tasks
- Depends on PostgreSQL and Redis

### Central Authentication
- Dockerfile: `Dockerfile.central-auth`
- Handles authentication for the system
- Depends on PostgreSQL and Redis

### PostgreSQL
- Uses PostgreSQL 16
- Exposes port 5433 (mapped to internal 5432)

### Redis
- Uses Redis 6
- Exposes port 6380 (mapped to internal 6379)

## Configuration

The application uses configuration file `config.yaml` located in the `config` directory. Ensure these are properly set up before running the services.

## Migrations

Database migrations are stored in the `migrations` directory. They are automatically applied when the services start up.

## Other instructions
1. To run unit test
    ```
    make test
    ```
2. To generate new migration file
    ```
    make create-migration name=<the-file-name> 
    ```