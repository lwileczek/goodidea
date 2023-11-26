# Good Idea
Good Idea is a slightly fancier TODO app where anonymous users may
submit requests and vote on their importance.
Users do **not** need to login or create an account to make a request or cast votes.

### Requirements
 - Postgresql Database
 - environment variable `DATABASE_URL` to point the server at the database

#### Dev Requirements
 - go 1.21+
 - node.js 18+

## Build
 1. Build and minify CSS and JS with: `npm install && npm run build`
 2. Build the server binary with `go build ./app/main.go`

 To make it easy, use [make](https://www.gnu.org/software/make/) to build

 ```bash
 make build
 ```

 which will populate the [static](./static) directory and build a binary called `server`

# Status
> **warning**
> This project is in _beta_

## Deployment
For deployment information, please see the project [wiki](https://github.com/lwileczek/goodidea/wiki)
