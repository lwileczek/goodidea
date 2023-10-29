# Good Idea
Write it down.
This is a simple tool to let users submit requests and vote on their importance.
Users to **not** need to login or create an account to make a requet or vote

### Requirements
 - Postgresql Database
 - `DATABASE_URL` to point to the database
 - The `templates` directory, I think

If there is no pre-built binary then go 1.21+ is required

## Build
 1. Minify any JS with esBuild: `npm install && npm run build`
 2. Build the binary with Go `go build ./...`

# Status
> **warning**
> This project is in _alpha_
