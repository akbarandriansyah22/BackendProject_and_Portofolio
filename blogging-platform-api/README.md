## Blogging Platform API

This project is an implementation of the Blogging Platform API challenge from roadmap.sh. It is a beginner-level backend project designed to demonstrate the fundamentals of building a RESTful service using Go and PostgreSQL.

## Features

- **CRUD Operations**: Complete functionality to Create, Read, Update, and Delete blog posts.

- **Search & Filtering**: Users can filter blog posts by search terms using query parameters.

- **PostgreSQL Integration**: Robust data storage using PostgreSQL for structured content management.

- **Layered Architecture**: The project follows a clean separation of concerns with dedicated layers for Handlers, Services, and Repositories.

## Technologies Used

- **Go**: The core programming language (v1.22+).

- **PostgreSQL**: The relational database used for data persistence.

- **github.com/lib/pq**: The Go driver for PostgreSQL communication.

- **Standard Library (net/http)**: Used for building the HTTP server and routing.

## Getting Started

To set up the project locally, follow these steps:

1. **Clone the repository**:

Bash

git clone https://github.com/yourusername/blogging-platform-api.git
cd blogging-platform-api

2. **Database Setup: Create a database named blogging_platform and execute the following SQL**:

SQL

CREATE TABLE posts (
id SERIAL PRIMARY KEY,
title VARCHAR(255) NOT NULL,
content TEXT NOT NULL,
category VARCHAR(100),
tags TEXT[],
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

3. **Configure Environment Variables: Set the following variables in your terminal or environment file**:

Bash

export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=blogging_platform

4. **Run the application**:

Bash

go mod download
go run cmd/main.go
The API server will be running at http://localhost:8080.

## API Endpoints

The following API endpoints are available:

**POST /posts**: Create a new blog post.

**GET /posts**: Retrieve all blog posts.

**GET /posts?q={keyword}**: Search posts by keywords in title or content.

**GET /posts/{id}**: Get details of a single post by ID.

**PUT /posts/{id}**: Update an existing post.

**DELETE /posts/{id}**: Remove a post from the database.
