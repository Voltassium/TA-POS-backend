# Agent Instructions — TA-POS-backend (backend-ta)

## Project Overview
This repository is a layered Go backend for a **Restaurant Point-of-Sale system** (module `backend-ta`). It provides a maintainable architecture with `internal/*` application code, `pkg/*` reusable packages, Bun-based DB access, MinIO-compatible object storage, JWT authentication, and Zap logging.

---

## Tech Stack

| Layer | Technology |
| :--- | :--- |
| Language | Go (>= 1.24) |
| Web Framework | Gin |
| Database | PostgreSQL |
| DB Library | Bun (github.com/uptrace/bun) |
| DB Driver | github.com/lib/pq |
| Auth | golang-jwt/jwt (v5) |
| Password Hashing | golang.org/x/crypto/bcrypt |
| Config | Viper (YAML files) |
| Logger | Zap (go.uber.org/zap) |
| Object Storage | MinIO (minio-go) |

---

## Project Layout

```
<repo-root>
├── main.go
├── go.mod
├── README.md
├── migrations/                         # Raw SQL migration files
│   ├── 000001_create_users.sql
│   ├── 000002_create_categories.sql
│   ├── 000003_create_products.sql
│   ├── 000004_create_orders.sql
│   ├── 000005_create_order_items.sql
│   └── 000006_create_payments.sql
├── domain/                             # Domain model types
│   ├── user.go
│   ├── category.go
│   ├── product.go
│   ├── order.go
│   ├── order_item.go
│   └── payment.go
├── dto/
│   ├── requests/
│   │   ├── auth_request.go
│   │   ├── user_request.go
│   │   ├── category_request.go
│   │   ├── product_request.go
│   │   ├── order_request.go
│   │   ├── order_item_request.go
│   │   └── payment_request.go
│   └── response/
│       ├── auth_response.go
│       ├── user_response.go
│       ├── category_response.go
│       ├── product_response.go
│       ├── order_response.go
│       ├── order_item_response.go
│       └── payment_response.go
├── repository/
│   ├── init.go
│   ├── user_repository.go
│   ├── category_repository.go
│   ├── product_repository.go
│   ├── order_repository.go
│   ├── order_item_repository.go
│   └── payment_repository.go
├── internal/
│   ├── controllers/
│   │   ├── auth_controller.go
│   │   ├── user_controller.go
│   │   ├── category_controller.go
│   │   ├── product_controller.go
│   │   ├── order_controller.go
│   │   ├── order_item_controller.go
│   │   └── payment_controller.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── category_service.go
│   │   ├── product_service.go
│   │   ├── order_service.go
│   │   ├── order_item_service.go
│   │   └── payment_service.go
│   └── routes/
│       └── routes.go
└── pkg/
    ├── authentication/
    ├── config/
    │   └── files/
    │       ├── env.yaml
    │       └── env.example.yaml
    ├── database/
    ├── http/
    │   └── server/
    │       └── http_response/
    ├── logger/
    ├── storage/
    └── validation/
```

---

## Configuration

Configuration uses Viper and lives at `pkg/config/files/env.yaml`. Copy from `env.example.yaml` to get started. **Never commit secrets to VCS.**

```yaml
database:
    host: localhost
    port: 5432
    user: postgres
    password: postgres
    name: restaurant_db
    ssl_mode: disable

application:
    port: 3000
    environment: development

authentication:
    encrypt_key: ""
    access_secret_key: ""
    refresh_token_key: ""
    access_token_expiry: "1h"
    refresh_token_expiry: "72h"
    issuer: "ta-pos-backend"

logger:
    environment: development
    log_level: debug
    encoding: json

object_storage:
    endpoint: ""
    access_key: ""
    secret_key: ""
```

---

## Database Schema

> **Important:** This project uses **explicit SQL migrations** — there is no AutoMigrate. Run all files in `migrations/` in order before starting the server.

### Migration Files

#### `000001_create_users.sql`
```sql
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    email       VARCHAR(255) UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    role        VARCHAR(50) NOT NULL DEFAULT 'staff', -- 'admin' or 'staff'
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### `000002_create_categories.sql`
```sql
CREATE TABLE IF NOT EXISTS categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    image_url   TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### `000003_create_products.sql`
```sql
CREATE TABLE IF NOT EXISTS products (
    id           SERIAL PRIMARY KEY,
    category_id  INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    price        NUMERIC(12, 2) NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### `000004_create_orders.sql`
```sql
CREATE TABLE IF NOT EXISTS orders (
    id           SERIAL PRIMARY KEY,
    table_id     INTEGER NOT NULL,
    staff_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
    status       VARCHAR(50) NOT NULL DEFAULT 'Open', -- 'Open', 'Paid', 'Cancelled'
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### `000005_create_order_items.sql`
```sql
CREATE TABLE IF NOT EXISTS order_items (
    id          SERIAL PRIMARY KEY,
    order_id    INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id  INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity    INTEGER NOT NULL CHECK (quantity > 0),
    unit_price  NUMERIC(12, 2) NOT NULL, -- snapshot of price at time of order
    subtotal    NUMERIC(12, 2) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### `000006_create_payments.sql`
```sql
CREATE TABLE IF NOT EXISTS payments (
    id              SERIAL PRIMARY KEY,
    order_id        INTEGER UNIQUE NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    payment_method  VARCHAR(50) NOT NULL, -- 'Cash', 'Card', 'Digital Wallet'
    amount_paid     NUMERIC(12, 2) NOT NULL,
    timestamp       TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## Domain Models

Define structs in `domain/` with Bun struct tags.

### `domain/user.go`
```go
type User struct {
    bun.BaseModel `bun:"table:users"`

    ID        int64     `bun:"id,pk,autoincrement"`
    Email     string    `bun:"email,notnull,unique"`
    Password  string    `bun:"password,notnull"`
    Role      string    `bun:"role,notnull,default:'staff'"`
    CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
```

### `domain/category.go`
```go
type Category struct {
    bun.BaseModel `bun:"table:categories"`

    ID        int64     `bun:"id,pk,autoincrement"`
    Name      string    `bun:"name,notnull"`
    ImageURL  string    `bun:"image_url"`
    Products  []Product `bun:"rel:has-many,join:id=category_id"`
    CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
```

### `domain/product.go`
```go
type Product struct {
    bun.BaseModel `bun:"table:products"`

    ID          int64     `bun:"id,pk,autoincrement"`
    CategoryID  int64     `bun:"category_id,notnull"`
    Name        string    `bun:"name,notnull"`
    Description string    `bun:"description"`
    Price       float64   `bun:"price,notnull"`
    IsAvailable bool      `bun:"is_available,notnull,default:true"`
    Category    *Category `bun:"rel:belongs-to,join:category_id=id"`
    CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
    UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
```

### `domain/order.go`
```go
type Order struct {
    bun.BaseModel `bun:"table:orders"`

    ID          int64       `bun:"id,pk,autoincrement"`
    TableID     int64       `bun:"table_id,notnull"`
    StaffID     int64       `bun:"staff_id,notnull"`
    TotalAmount float64     `bun:"total_amount,notnull,default:0"`
    Status      string      `bun:"status,notnull,default:'Open'"`
    OrderItems  []OrderItem `bun:"rel:has-many,join:id=order_id"`
    Payment     *Payment    `bun:"rel:has-one,join:id=order_id"`
    CreatedAt   time.Time   `bun:"created_at,nullzero,notnull,default:current_timestamp"`
    UpdatedAt   time.Time   `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
```

### `domain/order_item.go`
```go
type OrderItem struct {
    bun.BaseModel `bun:"table:order_items"`

    ID        int64     `bun:"id,pk,autoincrement"`
    OrderID   int64     `bun:"order_id,notnull"`
    ProductID int64     `bun:"product_id,notnull"`
    Quantity  int       `bun:"quantity,notnull"`
    UnitPrice float64   `bun:"unit_price,notnull"`
    Subtotal  float64   `bun:"subtotal,notnull"`
    Product   *Product  `bun:"rel:belongs-to,join:product_id=id"`
    CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
```

### `domain/payment.go`
```go
type Payment struct {
    bun.BaseModel `bun:"table:payments"`

    ID            int64     `bun:"id,pk,autoincrement"`
    OrderID       int64     `bun:"order_id,notnull,unique"`
    PaymentMethod string    `bun:"payment_method,notnull"`
    AmountPaid    float64   `bun:"amount_paid,notnull"`
    Timestamp     time.Time `bun:"timestamp,nullzero,notnull,default:current_timestamp"`
}
```

---

## API Endpoints

### Auth (Public)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| POST | `/v1/authentications/register` | Register a new user |
| POST | `/v1/authentications/login` | Login and receive JWT tokens |
| POST | `/v1/authentications/refresh-token` | Refresh access token |
| POST | `/v1/authentications/logout` | Logout (protected) |

### Users (Protected — Admin only)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| GET | `/v1/users` | Get all users |
| GET | `/v1/users/profile` | Get own profile |
| GET | `/v1/users/:id` | Get user by ID |
| PUT | `/v1/users/:id` | Update user |
| DELETE | `/v1/users/:id` | Delete user |

### Categories
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| GET | `/v1/categories` | Get all categories | Public |
| GET | `/v1/categories/:id` | Get category by ID | Public |
| POST | `/v1/categories` | Create category | Admin |
| PUT | `/v1/categories/:id` | Update category | Admin |
| DELETE | `/v1/categories/:id` | Delete category | Admin |

### Products
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| GET | `/v1/products` | Get all products | Public |
| GET | `/v1/products/:id` | Get product by ID | Public |
| POST | `/v1/products` | Create product | Admin |
| PUT | `/v1/products/:id` | Update product | Admin |
| DELETE | `/v1/products/:id` | Delete product | Admin |

### Orders
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| GET | `/v1/orders` | Get all orders | Staff + Admin |
| GET | `/v1/orders/:id` | Get order with items | Staff + Admin |
| POST | `/v1/orders` | Create a new order | Staff + Admin |
| PATCH | `/v1/orders/:id/status` | Update order status | Staff + Admin |
| DELETE | `/v1/orders/:id` | Cancel an order | Admin |

### Order Items
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| POST | `/v1/orders/:id/items` | Add item to order | Staff + Admin |
| DELETE | `/v1/orders/:id/items/:item_id` | Remove item from order | Staff + Admin |

### Payments
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| POST | `/v1/payments` | Process a payment | Staff + Admin |
| GET | `/v1/payments/:order_id` | Get payment by order | Staff + Admin |

---

## Business Logic Rules

### Orders
- New orders start with `status = "Open"`.
- `total_amount` must be **recalculated automatically** whenever an item is added or removed.
- Orders with `status = "Paid"` or `"Cancelled"` **cannot be modified**.

### Order Items
- When an item is added, copy `product.price` into `order_item.unit_price` as a **price snapshot**.
- Calculate `subtotal = unit_price * quantity` automatically.
- After any item change, recalculate and persist `order.total_amount`.
- Use `pkg/database.RunInTx` to wrap item insertion and order total update in a single transaction.

### Payments
- Payment can only be created if `order.status = "Open"`.
- After a successful payment, automatically set `order.status = "Paid"`.
- `amount_paid` must be **≥ order.total_amount**.
- Wrap payment creation and order status update in a single transaction.

### Auth & Roles
- Passwords must be hashed with **bcrypt** before saving.
- JWT is returned on login; validated on every protected route via middleware.
- Role `"admin"` has full access to all endpoints.
- Role `"staff"` can manage orders and payments but **cannot** manage users, categories, or products.

---

## Middleware

### `AuthMiddleware`
- Extract and validate JWT from the `Authorization: Bearer <token>` header.
- Attach user `id` and `role` to the Gin context.
- Return `401` if the token is missing or invalid.
- Lives in `pkg/http/server/middlewares`.

### `AdminOnly`
- Check if the role in context is `"admin"`.
- Return `403` if not.
- Lives in `pkg/http/server/middlewares`.

---

## Response Format

All responses use `http_response.SendSuccess` and `http_response.SendError` helpers.

### Success
```json
{
    "status": 200,
    "message": "OK",
    "data": {},
    "error": null
}
```

### Error
```json
{
    "status": 400,
    "message": "Bad Request",
    "data": null,
    "error": "Descriptive error message here"
}
```

---

## Startup Sequence (`main.go`)
1. Load configuration via `pkg/config.LoadConfig()`.
2. Initialize Zap logger via `pkg/logger.NewZapLogger(cfg.Logger)`.
3. Initialize Bun DB via `pkg/database.InitDB(cfg.Database)`.
4. Initialize JWT manager via `authentication.NewJWTManager(...)` and `authentication.SetupKey(...)`.
5. Initialize repositories via `repository.Init(database.GetDB())`.
6. Register all routes via `internal/routes.RegisterV1`.
7. Start HTTP server via `pkg/http/server.Init(cfg.Application, routes.RegisterV1)` with graceful shutdown.

> **Note:** Run all SQL migrations in `migrations/` manually before starting the server.

---

## How to Add a New Endpoint

1. **DTO** — Add request and response structs in `dto/requests/` and `dto/response/`.
2. **Repository** — Add DB query methods in `repository/`.
3. **Service** — Implement business logic in `internal/services/`.
4. **Controller** — Add the handler in `internal/controllers/`.
5. **Route** — Register the route in `internal/routes/routes.go`.

---

## Development Notes
- Use `pkg/database.RunInTx` for any multi-step DB operation (order creation, payment processing).
- Use `pkg/validation` helpers to validate all incoming request bodies.
- Never expose the `password` field in any response.
- Return proper HTTP status codes: `200`, `201`, `400`, `401`, `403`, `404`, `500`.
- Keep all secrets out of VCS; populate `env.yaml` from `env.example.yaml`.
- Use Zap logger for all logging — no `fmt.Println` in production code.
- Image uploads for categories should use the MinIO storage client in `pkg/storage`.