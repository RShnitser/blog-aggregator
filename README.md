# Gator

Gator is a multi-user CLI application. There's no server (other than the database), so it's only intended for local use.

## Installation
 
Go and Postgres database will need to be installed

```bash
go install ...
```

## Config

Create a config file in your home directory, ~/.gatorconfig.json, with the following content:
```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the values with your database connection string.

## Commands

Create a new user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start the aggregator:

```bash
gator agg 30s
```

View the posts:

```bash
gator browse [limit]
```

There are a few other commands you'll need as well:

- `gator login <name>` - Log in as a user that already exists
- `gator users` - List all users
- `gator feeds` - List all feeds
- `gator follow <url>` - Follow a feed that already exists in the database
- `gator unfollow <url>` - Unfollow a feed that already exists in the database