# Forum Web

A web forum built with Go and SQLite. Users can register, create posts with categories and images, comment, and like or dislike content.

## Prerequisites

- Go 1.21 or higher
- GCC (required by go-sqlite3 for CGO compilation)
  - Linux : `sudo apt install gcc`
  - macOS : included with Xcode Command Line Tools (`xcode-select --install`)
  - Windows : install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)

## Installation

```bash
git clone https://github.com/nito64chevrin-oss/Forum_web.git
cd Forum_web
go mod download
```

## Run

```bash
go run ./main
```

The server starts at [http://localhost:8080](http://localhost:8080).

## Run with Docker

```bash
docker build -t forum-web .
docker run -p 8080:8080 forum-web
```

## Features

- Register, login and logout
- Session management via cookies
- Create posts with categories and images
- Comment on posts
- Like and dislike posts and comments
- Filter posts by category, by created posts, or by liked posts
- Edit and delete your own posts and comments
- Edit your user profile
- Protected routes redirect unauthenticated users
- Custom 404 page

## Project Structure

```
Forum_web/
├── main/        # Go source code (handlers, database, routing)
├── static/      # CSS, JavaScript and uploaded images
├── views/       # HTML templates
└── forum.db     # SQLite database
```

## Dependencies

- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver
- [go.uuid](https://github.com/satori/go.uuid) - UUID generation for sessions
- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) - Password hashing (bcrypt)
