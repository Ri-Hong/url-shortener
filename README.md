

Add the following to your `.zshrc` file:

```bash
source <path-to-repo>/url-shortener/scripts/.zshrc
```

### Database

```bash
# To force a rebuild
docker-compose -f docker-compose-postgres.yml up --build

# To start the services in detached mode
docker compose -f docker-compose-postgres.yml up

# To stop the services
docker compose -f docker-compose-postgres.yml down

# To connect to the database
PGPASSWORD="mypassword" psql -h localhost -U myuser -d urlshortener

# To create a migration
migrate create -ext=sql -dir=infra/database/migrations -seq <name>
```
