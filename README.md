# resty-service

This is REST backend service.


## API Reference

#### health, ready, ping

```http request
    GET /healthz
    GET /readyz
    GET /ping
```

| Status |
|:-------|
| `200`  |



#### Get all items

```http request
  GET /api/items
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `api_key` | `string` | **Required**. Your API key |

#### Get item

```http request
  GET /api/items/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### add(num1, num2)

Takes two numbers and returns the sum.


## Deployment

To deploy this project run

```bash
  
```


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

```dotenv
SERVER_DB_URL=localhst:5432

SERVER_REDIS_URL=localhst:1111

SERVER_RABBIT_MQ_URL=localhst:11111

SERVER_HOST=localhst
SERVER_PORT=8080
SERVER_ENVIRONMENT=development (or production)
SERVER_WRITETIMEOUT=10s
SERVER_READTIMEOUT=10s
SERVER_IDLETIMEOUT=10s
```

#### Development

For development, if you are using IntelJ, `.run` folder will assist you with predefined environment variables.
You can them change whenever you want by clicking on `Edit configurations...`, selecting specific configuration and
clicking on the table with environment variables.

