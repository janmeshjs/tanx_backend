# Crypto Price Alert System

This project implements a Crypto Price Alert System that allows users to set alerts for specific cryptocurrencies and receive notifications when the target prices are reached.

## Setup and Run

### Prerequisites
- Docker Desktop installed on your machine.

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/crypto-price-alert.git
   cd crypto-price-alert
2. Build and run the Docker containers:
    ```bash
    docker-compose up --build
3. Access the application at http://localhost:8080.

## Endpoints

### Create Alert

Endpoint: POST /alerts/create
Description: Create a new price alert for a cryptocurrency.
Request Body:

  ```bash
  {
  "coin_id": "bitcoin",
  "target_price": 40000,
  "user_id": 123
  }
```
### Delete Alert
Endpoint: DELETE /alerts/delete/{id}
Description: Delete a price alert by its ID.

### Get Alerts
Endpoint: GET /alerts
Query Parameters:
page (optional): Page number for pagination.
status (optional): Filter alerts by status (created, deleted, triggered, etc.).

