# storex

# Storex – Asset Management System (Backend)

Storex is a **backend Asset Management System** built using **Golang** and **PostgreSQL**, designed with clean architecture principles, secure authentication, and role-based access control. This project was developed as part of a real-world backend assignment and focuses on scalability, clarity, and industry best practices.

---

## 🚀 Features

* Asset management (create, update, assign, track)
* User authentication using **JWT (JSON Web Tokens)**
* Role-Based Access Control (RBAC)

  * Admin
  * Sub-Admin
  * User
* Secure password hashing
* Middleware-based authentication & authorization
* PostgreSQL-backed persistent storage
* Clean, modular project structure
* Database migrations support

---

## 🛠 Tech Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Router:** Chi
* **Authentication:** JWT
* **Password Security:** bcrypt
* **Database Driver:** database/sql
* **Migration Tool:** golang-migrate
  
---

## 🔐 Authentication & Authorization

### Authentication

* JWT-based authentication
* Tokens are generated on successful login
* Tokens must be passed in the `Authorization` header:

```http
Authorization: Bearer <jwt_token>
```

### Authorization (RBAC)

Access is controlled using middleware based on user roles:

* **Admin:** Full access
* **asset_manager:** access only for asset management
* **employee_manager:** access only for employee management
* **employee:** only access to his own dashboard

---

## 🗄 Database Design

The database schema is designed incrementally with clarity and normalization in mind.

Core tables include:

* `users`
* `asset_brands`
* `asset_models`
* `assets`
* `services`
* `asset_status`
* `user_roles`
* `specs` table

All schema changes are managed via SQL migrations.

---

## ⚙️ Environment Variables

Create a `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=storex
```

---

## ▶️ Running the Project

### 1️⃣ Clone the repository

```bash
git clone <repository-url>
cd storex
```

### 3️⃣ Run the Server

```bash
go run cmd/main.go
```

Server will start on:

```
http://localhost:8080
```

---

## 🧪 API Design Philosophy

* RESTful endpoints
* Clear separation of concerns
* Thin handlers, fat services
* Middleware for cross-cutting concerns
* Explicit error handling

---

## 📌 Learning Outcomes

This project helped reinforce:

* Real-world backend architecture in Go
* JWT authentication internals
* Middleware design patterns
* Secure API development
* PostgreSQL schema planning
* Role-based authorization

---

## 👤 Author

**Gulshan Kumar**
Backend Developer (Golang) | AI/ML Enthusiast

---

## 📄 License

This project is for learning and demonstration purposes.
