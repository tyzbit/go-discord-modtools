# go-discord-modtools

Discord bot for Moder8or, written in Go

## Configuration

Set some environment variables before launching, or add a `.env` file.

If database environment variables are provided, the bot will save stats to an external database.
Otherwise, it will save stats to a local sqlite database at `/var/go-discord-modtools/local.db`

| Variable          | Value(s)                                                           |
| :---------------- | :----------------------------------------------------------------- |
| ADMINISTRATOR_IDS | Comma separated IDs of users allowed to use administrator commands |
| DB_NAME           | Database name for database                                         |
| DB_HOST           | Hostname for database                                              |
| DB_PASSWORD       | Password for database user                                         |
| DB_USER           | Username for database user                                         |
| LOG_LEVEL         | `trace`, `debug`, `info`, `warn`, `error`                          |
| TOKEN             | The Discord token the bot should use                               |

## Usage

Right click (Desktop) or long press (mobile) a message and select Apps, then select a moderation action for the message.

### Commands

Configure the bot:

`/settings`

Get reputation for a user:

`/reputation`

Get stats for the bot:

`/stats`

Get this help message:

`/help`

## Development

Create a `.env` file with your configuration, at the bare minimum you need
a Discord token for `TOKEN` and an Moder8or Cookie for `COOKIE` (Need at least `PHPSESSID`, `logged-in-sig` and `logged-in-user`, looks like: `PHPSESSID=12345; logged-in-sig=54321; logged-in-user=example%40example.com`).

Logins are currently good for a year.
You can either `docker compose up --build` to run with a mysql database, or just `go run main.go` to run with a sqlite database.
