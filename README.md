# MenuLink-up

A full stack application written in Go and Svelte for uploading and managing PDFs using S3

![demo](https://github.com/LouisHatton/menu-link-up/assets/71732103/fff85436-9f73-4070-911d-86ded046c891)

## Running in development

Example env files are provided for the api and ui at `.env.example` and `cmd/ui/.env.example` respectively.

Start the docker compose file which creates the MariaDB container. 
The API should automatically add any required tables.

```
$ docker compose up -d
```

Run the API with the command:

```
$ go run ./cmd/api
```

Finally run the UI with the command:

```
$ cd ./cmd/ui && npm run dev
```
