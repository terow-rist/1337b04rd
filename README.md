# 1337b04rd

An anonymous imageboard backend written in Go, featuring thread and comment posting, image uploads to S3-compatible storage, and session tracking via cookies. Built with Hexagonal Architecture and clean coding principles, it includes a frontend rendered with Go templates.


## ✨ Features

* Anonymous sessions with Rick & Morty avatars
* Thread and comment posting
* Auto-cleanup of threads without comments
* Moderation-ready architecture
* Clean Go codebase with layered separation

---

## 🛠️ Tech Stack

* Go 1.23+
* PostgreSQL
* HTML templates (server-rendered)
* Docker (virtualization)

---

## 📁 Project Structure

```bash
1337b04rd/
.
├── cmd
│   └── app.go
├── compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── adapter
│   │   ├── api
│   │   │   └── rick_and_morty.go
│   │   ├── config
│   │   │   ├── config.go
│   │   │   └── flag.go
│   │   ├── handler
│   │   │   └── http
│   │   │       ├── helpers.go
│   │   │       ├── middleware.go
│   │   │       ├── posts.go
│   │   │       └── router.go
│   │   ├── logger
│   │   │   └── slog.go
│   │   └── storage
│   │       ├── postgres
│   │       │   ├── db.go
│   │       │   ├── migrations
│   │       │   │   ├── 001_create_user.sql
│   │       │   │   ├── 002_create_posts.sql
│   │       │   │   └── 003_create_comments.sql
│   │       │   └── repository
│   │       │       ├── comments.go
│   │       │       ├── posts.go
│   │       │       └── users.go
│   │       └── s3
│   └── core
│       ├── domain
│       │   ├── comment.go
│       │   ├── post.go
│       │   └── user.go
│       ├── port
│       │   ├── avatar.go
│       │   ├── comments.go
│       │   ├── error.go
│       │   ├── posts.go
│       │   └── users.go
│       ├── service
│       │   ├── comments.go
│       │   ├── posts.go
│       │   └── users.go
│       └── util
├── main.go
├── Makefile
├── README.md
└── templates
    ├── archive.html
    ├── archive-post.html
    ├── catalog.html
    ├── create-post.html
    ├── error.html
    └── post.html
```



---

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/fallen-fatalist/1337b04rd.git
cd 1337b04rd
```



### 2. Set Up Environment Variables

Create a `.env` file based on `.env.example` and configure the following variables:([GitHub][2])

* **HTTP Server**

  * `PORT=8080`
* **PostgreSQL**

  * `DB_HOST=localhost`
  * `DB_PORT=5432`
  * `DB_USER=your_db_user`
  * `DB_PASSWORD=your_db_password`
  * `DB_NAME=1337b04rd`
  * `DB_SSLMODE=disable`


### 3. Run the Application

Use Docker Compose to build and run the application:

```bash
docker-compose up --build
```



---

## 💻 Usage

Then, open your browser and go to `http://localhost:8000/` 

---

## ⚙️ Configuration

All configuration settings are managed through the `.env` file. Ensure all required environment variables are set before running the application.

---

## 🧪 Testing

Run the test suite using the following command:

```bash
go test -v ./...
```


