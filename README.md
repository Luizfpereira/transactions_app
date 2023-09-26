# Transactions App

Transactions app is an app that allows users to store purchases transactions in a database and retrieve them with a specific currency and exchange rate.

## Functionalities

- Transactions creation
- Query transactions by currency
- Query a transaction by ID and currency

## Prerequisites

- Docker
- Docker-compose

## Installation

Clone this repository:
```shell
git clone https://github.com/Luizfpereira/transactions_app.git
```

## Running Transactions App
1. Navigate to the application repository:
```shell
cd transactions_app
```

2. Inside the repository in root, execute the command:
```shell
docker-compose up -d
```

The database and the application will be running on ports 5432 and 8080, respectively

## Usage

### Routes


Method | Endpoint | Description
------ |--------- | -----------
GET    | /        | Check if application is running
POST   | /transactions/create | Create a new transaction
GET   | /transactions | Query transactions by currency
GET | /transactions/{id} | Query a transaction by ID and currency

## Example usage

### POST /transactions/create

To create a new transaction, make a POST request to `/transactions/create` with the transaction data in the request body as form data

The transaction date should be in the format `YYYY-MM-DD HH:MM:SS`. Example: `2021-09-20 22:00:00`

Example:
```shell
curl -X POST http://localhost:8080/transactions/create -F description='test' -F transaction_date='2021-09-20 22:00:00' -F purchase_amount=959.46
```

Response:

Status code: 200 OK
```json
{
    "message": {
        "ID": 1,
        "CreatedAt": "2023-09-26T01:48:46.181898895Z",
        "UpdatedAt": "2023-09-26T01:48:46.181898895Z",
        "DeletedAt": null,
        "TransactionDate": "2021-09-20T22:00:00Z",
        "Description": "test",
        "PurchaseAmount": "959.46"
    },
    "status": "success"
}
```

Empty description:

Status code: 400 Bad Request
```json
{
    "error": "description is empty",
    "status": "failed"
}
```

Description length greater than 50 characters:

Status code: 400 Bad Request
```json
{
    "error": "description max length: 50",
    "status": "failed"
}
```

Empty transaction date or malformed:

Status code: 400 Bad Request
```json
{
    "error": "failed parsing date: parsing time \"\" as \"2006-01-02 15:04:05\": cannot parse \"\" as \"2006\"",
    "status": "failed"
}
```

---

### GET /transactions

This endpoint queries stored purchase data with or without specifying a currency. The currencies available can be found in the documentation of the Treausy Fiscal Data: `https://fiscaldata.treasury.gov/api-documentation/#getting-started`

To specify a currency, one should use a query param `?currency={currency name}`. Example: `localhost:8080/transactions?currency=Brazil-Real`

#### No currency: /transactions

Example:
```shell
curl http://localhost:8080/transactions
```

Response:

Status code: 200 OK
```json
{
    "message": [
        {
            "id": 1, 
            "description": "test1",
            "transaction_date": "2023-09-20T22:00:00Z",
            "purchase_amount": "506.7"
        },
        {
            "id": 2, 
            "description": "test2",
            "transaction_date": "2023-09-20T22:00:00Z",
            "purchase_amount": "1000"
        },
        {
            "id": 3, 
            "description": "test3",
            "transaction_date": "2021-09-20T22:00:00Z",
            "purchase_amount": "1000",
        },
    ],
    "status": "success"
}
```

#### With currency: /transactions?currency={Country-Currency_name}

Example:
```shell
curl http://localhost:8080/transactions?currency=Brazil-Real
```

Response:

Status code: 200 OK
```json
{
    "message": [
        {
            "id": 1, 
            "description": "test1",
            "transaction_date": "2023-09-20T22:00:00Z",
            "purchase_amount": "506.7",
            "exchange_rate": "4.858",
            "converted_amount": "2461.55"
        },
        {
            "id": 2, 
            "description": "test2",
            "transaction_date": "2015-09-20T22:00:00Z",
            "purchase_amount": "1000",
            "error": "purchase cannot be converted to target currency"
        },
        {
            "id": 3, 
            "description": "test3",
            "transaction_date": "2021-09-20T22:00:00Z",
            "purchase_amount": "1000",
            "exchange_rate": "4.96",
            "converted_amount": "4960"
        },
    ],
    "status": "success"
}
```

---
### GET /transactions/{id}

This endpoint queries stored purchase data by id with or without specifying a currency. The currencies available can be found in the documentation of the Treausy Fiscal Data: `https://fiscaldata.treasury.gov/api-documentation/#getting-started`

To specify a currency, one should use a query param `?currency={currency name}`. Example: `localhost:8080/transactions?currency=Brazil-Real`

#### No currency: /transactions

Example:
```shell
curl http://localhost:8080/transactions/2
```

Response:

Status code: 200 OK
```json
{
    "message": {
        "id": 2,
        "description": "teste",
        "transaction_date": "2023-09-20T22:00:00Z",
        "purchase_amount": "506.7"
    },
    "status": "success"
}
```

#### With currency: /transactions?currency={Country-Currency_name}

Example:
```shell
curl http://localhost:8080/transactions/2?currency=Brazil-Real
```

Response:

Status code: 200 OK
```json
{
    "message": {
        "id": 2,
        "description": "teste",
        "transaction_date": "2023-09-20T22:00:00Z",
        "purchase_amount": "506.7",
        "exchange_rate": "4.858",
        "converted_amount": "2461.55"
    },
    "status": "success"
}
```

ID not registered:

Status code: 400 Bad Request
```json
{
    "error": "id not registered",
    "status": "failed"
}
```












