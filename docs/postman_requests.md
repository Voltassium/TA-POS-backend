# Postman Requests — TA-POS-backend (v1)

Base environment

- `base_url` = `http://localhost:8080/v1`
- `access_token` (set after login)
- `refresh_token` (set after login)

Headers to use for JSON requests

- `Content-Type: application/json`
- For protected endpoints: `Authorization: Bearer {{access_token}}`

---

## Quick notes

- All responses are wrapped using the server response envelope:

```json
{
  "status": <http-status-code>,
  "message": "<text>",
  "data": <object|null|array>,
  "error": null
}
```

Use `{{base_url}}` as the base for every request.

---

## Ping

- Method: `GET`
- URL: `{{base_url}}/ping`
- Auth: no
- Body: none

Example success (200):

```json
{
  "status": 200,
  "message": "pong",
  "data": {
    "time": "2026-04-20T12:34:56Z",
    "ua": "curl/7.80.0"
  },
  "error": null
}
```

---

## Authentication

### Register

- Method: `POST`
- URL: `{{base_url}}/authentications/register`
- Auth: no
- Headers: `Content-Type: application/json`
- Body (JSON):

```json
{
  "email": "admin@example.com",
  "password": "Secret123",
  "role": "admin"
}
```

Success (201):

```json
{
  "status": 201,
  "message": "Created",
  "data": null,
  "error": null
}
```

### Login

- Method: `POST`
- URL: `{{base_url}}/authentications/login`
- Auth: no
- Body (JSON):

```json
{
  "email": "admin@example.com",
  "password": "Secret123"
}
```

Success (200):

```json
{
  "status": 200,
  "message": "Login successful",
  "data": {
    "access_token": "<jwt_access_token>",
    "refresh_token": "<jwt_refresh_token>"
  },
  "error": null
}
```

Save the `access_token`/`refresh_token` into Postman environment variables after login.

### Refresh token

- Method: `POST`
- URL: `{{base_url}}/authentications/refresh-token`
- Auth: no
- Body (JSON):

```json
{ "refresh_token": "{{refresh_token}}" }
```

Success (200):

```json
{
  "status": 200,
  "message": "Refresh token successful",
  "data": { "access_token": "<new_access_token>" },
  "error": null
}
```

### Logout

- Method: `POST`
- URL: `{{base_url}}/authentications/logout`
- Auth: yes (Bearer)
- Headers: `Authorization: Bearer {{access_token}}`
- Body: none

Success (200):

```json
{
  "status": 200,
  "message": "Logout successful",
  "data": null,
  "error": null
}
```

---

## Users (admin role required)

> All `/users` endpoints require `Authorization: Bearer {{access_token}}` and admin role.

### List users

- Method: `GET`
- URL: `{{base_url}}/users`
- Query params (optional): `page`, `page_size`, `order_by`, `order_dir`

Example success (200):

```json
{
  "status": 200,
  "message": "Success get list data",
  "data": {
    "current_page": 1,
    "page_size": 10,
    "total_items": 1,
    "total_pages": 1,
    "has_previous": false,
    "has_next": false,
    "data": [
      {"id":1,"email":"admin@example.com","role":"admin","created_at":"2026-04-20T12:00:00Z","updated_at":"2026-04-20T12:00:00Z"}
    ]
  },
  "error": null
}
```

### Get profile

- Method: `GET`
- URL: `{{base_url}}/users/profile`
- Auth: yes

Success (200): returns a single user object in `data`.

### Get user

- Method: `GET`
- URL: `{{base_url}}/users/:id`
- Auth: yes (admin)

### Update user

- Method: `PUT`
- URL: `{{base_url}}/users/:id`
- Auth: yes (admin)
- Body example (any of the fields):

```json
{ "email": "new@example.com", "password": "NewSecret", "role": "staff" }
```

Success (200):

```json
{ "status": 200, "message": "OK", "data": null, "error": null }
```

### Delete user

- Method: `DELETE`
- URL: `{{base_url}}/users/:id`
- Auth: yes (admin)

Success (200): same OK envelope.

---

## Categories

### List categories (public)

- Method: `GET`
- URL: `{{base_url}}/categories`
- Auth: no

Success: pagination result containing category objects:

```json
{
  "status": 200,
  "message": "Success get list data",
  "data": {
    "current_page": 1,
    "page_size": 10,
    "total_items": 1,
    "total_pages": 1,
    "has_previous": false,
    "has_next": false,
    "data": [{"id":1,"name":"Drinks","image_url":"https://...","created_at":"2026-04-20T...","updated_at":"2026-04-20T..."}]
  },
  "error": null
}
```

### Get category (public)

- Method: `GET`
- URL: `{{base_url}}/categories/:id`

### Create category (admin)

- Method: `POST`
- URL: `{{base_url}}/categories`
- Auth: yes (admin)
- Body example:

```json
{ "name": "Drinks", "image_url": "https://example.com/drinks.jpg" }
```

Success (201):

```json
{
  "status": 201,
  "message": "Created",
  "data": {"id":10,"name":"Drinks","image_url":"https://example.com/drinks.jpg","created_at":"2026-04-20T...","updated_at":"2026-04-20T..."},
  "error": null
}
```

### Update/Delete category (admin)

- `PUT {{base_url}}/categories/:id` body similar to create (fields optional)
- `DELETE {{base_url}}/categories/:id` — both return `200 OK` envelope on success.

---

## Products

### List products (public)

- Method: `GET`
- URL: `{{base_url}}/products`
- Query params: `category_id`, `page`, `page_size`

### Get product (public)

- Method: `GET`
- URL: `{{base_url}}/products/:id`

### Create product (admin)

- Method: `POST`
- URL: `{{base_url}}/products`
- Auth: yes (admin)
- Body example:

```json
{
  "category_id": 1,
  "name": "Burger",
  "description": "Tasty",
  "price": 5.99,
  "is_available": true
}
```

Success (201): returns created product object in `data`.

### Update/Delete product (admin)

- `PUT {{base_url}}/products/:id` (optional fields)
- `DELETE {{base_url}}/products/:id`

Both return `200 OK` envelope on success.

---

## Orders (roles: admin, staff)

### List orders

- Method: `GET`
- URL: `{{base_url}}/orders`
- Query params: `page`, `page_size`, `status` (valid: `Open`,`Paid`,`Cancelled`)

### Get order

- Method: `GET`
- URL: `{{base_url}}/orders/:id`

Success (200) returns order detail:

```json
{
  "status": 200,
  "message": "Success get data",
  "data": {
    "id": 123,
    "table_id": 12,
    "staff_id": 2,
    "total_amount": 25.5,
    "status": "Open",
    "created_at": "2026-04-20T...",
    "updated_at": "2026-04-20T...",
    "items": [
      {"id":1,"order_id":123,"product_id":10,"product_name":"Burger","quantity":2,"unit_price":5.5,"subtotal":11.0,"created_at":"...","updated_at":"..."}
    ],
    "payment": null
  },
  "error": null
}
```

### Create order

- Method: `POST`
- URL: `{{base_url}}/orders`
- Body example:

```json
{ "table_id": 12 }
```

Success (201): created order object in `data`.

### Update order status

- Method: `PATCH`
- URL: `{{base_url}}/orders/:id/status`
- Body example:

```json
{ "status": "Paid" }
```

Success (200): `OK` envelope.

### Cancel order

- Method: `DELETE`
- URL: `{{base_url}}/orders/:id`

Success (200): `OK` envelope.

### Add order item

- Method: `POST`
- URL: `{{base_url}}/orders/:id/items`
- Body example:

```json
{ "product_id": 10, "quantity": 2 }
```

Success (200) returns the added item or updated items in `data`.

### Remove order item

- Method: `DELETE`
- URL: `{{base_url}}/orders/:id/items/:item_id`

Success (200) returns updated items or confirmation in `data`.

---

## Payments (roles: admin, staff)

### Create payment

- Method: `POST`
- URL: `{{base_url}}/payments`
- Auth: yes (admin|staff)
- Body example:

```json
{
  "order_id": 123,
  "payment_method": "Cash",
  "amount_paid": 50.0
}
```

Success (201):

```json
{
  "status": 201,
  "message": "Created",
  "data": {"id":1,"order_id":123,"payment_method":"Cash","amount_paid":50.0,"timestamp":"2026-04-20T..."},
  "error": null
}
```

### Get payment by order

- Method: `GET`
- URL: `{{base_url}}/payments/:order_id`

Success (200): returns payment object in `data`.

---

## Tips for Postman

- Create an environment with `base_url`, `access_token`, `refresh_token`.
- Use the Login request to set `access_token`/`refresh_token` and add `Authorization` header.
