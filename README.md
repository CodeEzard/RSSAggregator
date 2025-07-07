# RSSAggregator

A modern, concurrent RSS feed aggregator written in Go. Collects posts from multiple RSS feeds, cleans and stores them in PostgreSQL, and exposes a REST API for users to register, follow feeds, and retrieve posts.

---

## Features

- 🚀 Fast concurrent RSS scraping (configurable interval and workers)
- 🧹 Automatic HTML and UTF-8 cleaning for all post data
- 🔑 User registration and API key authentication
- 📡 Add, follow, and unfollow RSS feeds
- 💼 Aggregates job posts and other content
- 🗄️ PostgreSQL for persistent storage
- 🌐 REST API with pretty JSON responses

---

## Project Structure

```
RSSAggregator/
├── cmd/api/                 # Main application entry point
├── internal/
│   ├── handlers/            # HTTP handlers (users, feeds, posts, etc.)
│   ├── database/            # Database queries and models (sqlc generated)
│   ├── scraper/             # RSS scraping logic
│   ├── utils/               # Utilities (JSON, cleaning, etc.)
│   ├── models/              # Data models and transformations
│   └── auth/                # Authentication helpers
├── sql/
│   ├── schema/              # Database migrations
│   └── queries/             # SQL queries for sqlc
├── .env                     # Environment variables
├── go.mod / go.sum          # Go dependencies
└── README.md
```

---

## Quick Start

### 1. Clone & Install

```bash
git clone https://github.com/CodeEzard/RSSAggregator.git
cd RSSAggregator
go mod tidy
```

### 2. Database Setup

```bash
# Create the database (if not exists)
createdb rssaggregator

# Run migrations
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/rssaggregator?sslmode=disable" up
```

### 3. Environment Variables

Create a `.env` file:

```
PORT=8080
DB_URL=postgres://username:password@localhost:5432/rssaggregator?sslmode=disable
```

### 4. Generate Database Code

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc generate
```

### 5. Build & Run

```bash
go build -o bin/rssaggregator cmd/api/main.go
./bin/rssaggregator
```

---

## API Overview

All endpoints (except user creation) require:

```
Authorization: ApiKey YOUR_API_KEY
```

### User

- `POST /v1/users` — Register a new user
- `GET /v1/users` — Get current user info

### Feeds

- `POST /v1/feeds` — Add a new RSS feed
- `GET /v1/feeds` — List all feeds

### Feed Follows

- `POST /v1/feed_follows` — Follow a feed
- `GET /v1/feed_follows` — List followed feeds
- `DELETE /v1/feed_follows/{feedFollowID}` — Unfollow a feed

### Posts

- `GET /v1/posts` — Get posts from followed feeds
- `GET /v1/posts?clean=true&pretty=true` — Get cleaned, pretty-printed posts

---

## Example Usage

```bash
# Register a user
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'

# Add a feed
curl -X POST http://localhost:8080/v1/feeds \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "RemoteOK", "url": "https://remoteok.io/remote-jobs.rss"}'

# Follow a feed
curl -X POST http://localhost:8080/v1/feed_follows \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"feed_id": "FEED_UUID"}'

# Get posts
curl -H "Authorization: ApiKey YOUR_API_KEY" \
  "http://localhost:8080/v1/posts?pretty=true"
```

---

## Data Cleaning

- Removes HTML tags and entities
- Cleans invalid UTF-8 sequences
- Normalizes whitespace
- Truncates long descriptions

---

## Troubleshooting

- **UTF-8 errors**: The app automatically cleans invalid byte sequences.
- **Database errors**: Check your `DB_URL` and PostgreSQL status.
- **Feed errors**: Ensure the RSS feed URL is valid and reachable.

---

## Contributing

1. Fork the repo
2. Create a feature branch
3. Commit your changes
4. Open a pull request

---

## License

MIT License

---

**Built with Go,