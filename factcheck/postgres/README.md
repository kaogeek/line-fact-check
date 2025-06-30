# Database Schema and SQLC Setup

This directory contains the database schema and SQLC configuration for the factcheck project.

## Files

- `schema.sql` - Database schema with tables for all structs in `factcheck.go`
- `query.sql` - SQL queries for CRUD operations
- `sqlc.yaml` - SQLC configuration file

## Tables

### Topics
Stores fact-checking topics with their status and results.

### Messages
Stores messages associated with topics.

### UserMessages
Stores user message metadata with generic JSONB field for flexible data storage.

## Usage

To generate Go code from SQL:

```bash
cd factcheck
sqlc generate
```

This will generate Go code in the `postgres/` directory with:
- Database models matching the structs in `factcheck.go`
- CRUD operations for all tables
- Type-safe database interface

## Relationships

- `messages.topic_id` → `topics.id` (CASCADE delete)
- `user_messages.message_id` → `messages.id` (CASCADE delete)

## Indexes

Performance indexes are created on:
- `topics.status` and `topics.created_at`
- `messages.topic_id` and `messages.created_at`
- `user_messages.message_id` and `user_messages.created_at` 