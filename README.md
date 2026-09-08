# Balaji Clothing - Full Stack E-Commerce Platform

A complete, production-style e-commerce website built with **Go (Gin)** backend, **PostgreSQL**, **GORM**, **JWT authentication**, and a modern responsive frontend (HTML/CSS/JS).

## Features

- **User Authentication**: Register, Login, JWT, bcrypt password hashing, role-based access (user/admin)
- **Product Catalog**: Browse, search, filter by category/price, sort, pagination, product details
- **Shopping Cart**: Add/remove/update quantities, stock validation, persistent cart for logged-in users
- **Checkout & Orders**: Full shipping form, COD + mock online payment, order generation, stock reduction
- **Order Management**: View order history, statuses (pending → confirmed → shipped → delivered / cancelled)
- **Admin Panel**: Dashboard, manage products (CRUD), manage orders (status updates), view users
- **Security**: JWT protected routes, admin authorization, input validation, GORM (SQL injection safe)
- **UI/UX**: Modern responsive design, loading states, empty states, toast notifications

## Tech Stack

| Layer     | Technology                          |
|-----------|-------------------------------------|
| Backend   | Go 1.22+, Gin framework             |
| Database  | PostgreSQL + GORM                   |
| Auth      | JWT (golang-jwt) + bcrypt           |
| Frontend  | HTML5, CSS3, Vanilla JavaScript     |
| CORS      | gin-contrib/cors                    |

## Project Structure

```
balaji-shop/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/                 # Environment configuration
│   ├── database/               # DB connection, migrate, seed
│   ├── handlers/               # HTTP handlers (auth, product, cart, order)
│   ├── middleware/             # JWT auth & admin middleware
│   └── models/                 # GORM models & request DTOs
├── routes/routes.go            # Route definitions
├── frontend/                   # HTML pages
├── static/css|js               # Assets
├── bin/server                  # Compiled binary
└── README.md
```

## Prerequisites

- Go 1.22+
- PostgreSQL 14+

## Setup

### 1. PostgreSQL

```bash
# Create user and database
sudo -u postgres psql -c "CREATE USER ecommerce WITH PASSWORD 'ecommerce123';"
sudo -u postgres psql -c "CREATE DATABASE ecommerce OWNER ecommerce;"
```

### 2. Environment Variables (optional)

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=ecommerce
export DB_PASSWORD=ecommerce123
export DB_NAME=ecommerce
export JWT_SECRET=your-secret-key
export PORT=8080
```

Defaults are already set for local development.

### 3. Install & Run

```bash
cd balaji-shop
go mod tidy
go build -o bin/server ./cmd/server
./bin/server
```

Server starts at **http://localhost:8080**

On first run the database is auto-migrated and seeded with sample data.

## Default Accounts

| Role  | Email              | Password  |
|-------|--------------------|-----------|
| Admin | admin@balaji.com   | admin123  |
| User  | user@balaji.com    | user123   |

## API Endpoints

### Auth
- `POST /api/auth/register` – Register
- `POST /api/auth/login` – Login
- `GET  /api/auth/profile` – Get profile (auth)
- `PUT  /api/auth/profile` – Update profile (auth)

### Products
- `GET /api/products` – List (supports search, category_id, min_price, max_price, sort_by, page, limit)
- `GET /api/products/featured` – Featured products
- `GET /api/products/:id` – Product detail
- `GET /api/categories` – List categories

### Cart (auth required)
- `GET    /api/cart`
- `POST   /api/cart/items`
- `PUT    /api/cart/items/:id`
- `DELETE /api/cart/items/:id`
- `DELETE /api/cart`

### Orders (auth required)
- `POST /api/checkout`
- `GET  /api/orders`
- `GET  /api/orders/:id`

### Admin (auth + admin role)
- `GET    /api/admin/products`
- `POST   /api/admin/products`
- `PUT    /api/admin/products/:id`
- `DELETE /api/admin/products/:id`
- `GET    /api/admin/orders`
- `PUT    /api/admin/orders/:id/status`
- `GET    /api/admin/users`

## Frontend Pages

| Path              | Description                |
|-------------------|----------------------------|
| `/`               | Home / Landing             |
| `/shop`           | Product listing + filters  |
| `/product?id=`    | Product detail             |
| `/cart`           | Shopping cart              |
| `/checkout`       | Checkout form              |
| `/login`          | Login                      |
| `/register`       | Registration               |
| `/orders`         | Order history              |
| `/profile`        | User profile               |
| `/admin`          | Admin dashboard            |
| `/admin/products` | Product management         |
| `/admin/orders`   | Order management           |

## Payment Architecture

Supports:
- **Cash on Delivery** (`cod`)
- **Mock Online Payment** (`mock_online`) – marks payment as paid and order as confirmed

The checkout handler is structured so real gateways (Razorpay / Stripe) can be plugged in later without changing the order flow.

## Testing the Flows

1. Register a new user or login with `user@balaji.com`
2. Browse / search / filter products
3. Add items to cart, change quantities
4. Checkout with shipping details (COD or mock online)
5. View order in My Orders
6. Login as admin → manage products & order statuses

## License

MIT – Built for educational / demonstration purposes.
