# Golang Flashcards API

This is an example project to show how to deploy a Go app with docker-compose.

## Tech Stack

- Go 1.24
- PostgreSQL 16+
- sqlx
- lib/pq
- standard library (net/http)

## Deployment

This is a classic app managed with docker-compose. You can deploy it on any server with docker engine installed.

### Deploy with Dokploy

This works with any Platform-as-a-Service (PaaS) provider, like Coolify, Dokploy, etc.

- Get a server (VPS/DS)
- Install docker engine
- Fork a repository
- Install Dokploy on your server
- Create a GitHub app with access to your forked repository
- Create a new project
- Add a new service connected to your GitHub repository that runs this app
- Configure the repository: domain, environment variables, etc.
- Add a database with the image from `docker-compose.yml`
- Deploy it, make it publicly available (temporarily)
- Configure .env locally on your machine to connect to the database
- Run `make migrate-up` to apply migrations
- With psql CLI, run the `scripts/init.sql` script to create the extension on the database
- Disable public access for DB for security unless it's a test environment
- Update the app's environment variables
- Start the app
- Use it
