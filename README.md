Card Games API
==============

Backend API to support the implementation of playing card games.

Currently implements a service to manage a deck of playing cards - standard 52 cards deck without Joker cards.

**Technology stack**

* [Gin](https://github.com/gin-gonic/gin) as REST/Web library.
* [Gorm](https://gorm.io/index.html) as database layer (ORM).
* Database: supports PostgreSQL and Sqlite.
* Provided `Dockerfile` with multistage build for containerization.
* Example deployment with `docker-compose`.

# Building and running
 
## Build and run locally with `go`

To build the api locally with Go, first make sure you have [Go installed](https://go.dev/doc/install).
Go 1.18 or an earlier version like 1.17 should work just fine.

Checkout the source code with `git`:
```bash
git clone https://github.com/natemago/card-games-api.git
```

Then cd into the source code directory and build it:
```bash
cd card-games-api
go build -o card-games-api
```

To run a simple local version using Sqlite, run:

```bash
./card-games-api --db-type="sqlite" --db-url="card-games.db"
```
This will run the service at port `8080`.

To connect to a PostgreSQL database or to bind at different host and port, please refer to the [configuration](#configuration) for a list of parameters or ENV variables.

## Build and run with `docker`

To build with docker, first make sure you have [Docker installed](https://docs.docker.com/get-docker/) on your machine.

Checkout the source code with `git`:
```bash
git clone https://github.com/natemago/card-games-api.git
```

Then cd into the source code directory and build a docker image:

```bash
cd card-games-api

# Build docker image
docker build -t card-games-api:latest .
```

Run the docker container:

```bash
docker run -d -p 8080:8080 -e "DB_TYPE=sqlite" -e "DB_URL=/database.db" card-games-api:latest
```

Confirm that the app is running with `docker ps`:
```
docker ps
CONTAINER ID   IMAGE                   COMMAND                  CREATED         STATUS         PORTS                                       NAMES
77783f97810b   card-games-api:latest   "./card-games-api"       2 seconds ago   Up 2 seconds   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   competent_brown

```

This will run the api app withing a docker container with port 8080 exposed on the host.

To connect to a PostgreSQL database or to bind at different host and port, please refer to the [configuration](#configuration) for a list of parameters or ENV variables.

## Build and run with `docker-compose`

To build and run with Docker Compose, first [install `docker-compose`](https://docs.docker.com/compose/install/).

Checkout the source code with `git` and cd into the source code directory:
```bash
git clone https://github.com/natemago/card-games-api.git

cd card-games-api
```

A `docker-compose.yaml` is provided that sets up the API with PostgreSQL database.

Run it:
```bash
docker-compose up
```

To stop it and clean up hit CTRL-C, then `docker-compose down`.

To rebuild the images:
```bash
docker-compose build
```

# Configuration
The API can be configured and run using either the provided program flags or via ENV variables.
## ENV variables

The following ENV variables are available for configuration:

* `DB_URL` - the URL or DSN of the database. This is the database connection string.
  * For PostgreSQL, you can supply the DSN, for example: 
  `DB_URL="host=postgres user=toggl_user password=toggl_password dbname=toggl_card_games port=5432"`

  * For sqlite, you can provide the file name: `DB_URL="my-database.db"`
* `DB_TYPE` - is the database type (a db dialect) to use. For PostgreSQL set this to `postgres`; for sqlite set this to `sqlite`.
* `BIND_HOST` - the hostname to bind to when starting the HTTP server. By default this is set to empty string `""` - basically bind to all interfaces.
* `BIND_PORT` - on which port to listen for incoming HTTP connections. The default port is `8080`.

## Start parameters

The configuration parameters can also be controlled with program arguments.
The following flags can be used:

* `--db-url` - the URL or DSN of the database. This is the database connection string.
  * For PostgreSQL, you can supply the DSN, for example: 
  `DB_URL="host=postgres user=toggl_user password=toggl_password dbname=toggl_card_games port=5432"`

  * For sqlite, you can provide the file name: `DB_URL="my-database.db"`
* `--db-type` - is the database type (a db dialect) to use. For PostgreSQL set this to `postgres`; for sqlite set this to `sqlite`. The default value is `postgres`.
* `--bind-host` - the hostname to bind to when starting the HTTP server. By default this is set to empty string `""` - basically bind to all interfaces.
* `--bind-port` - on which port to listen for incoming HTTP connections. The default port is `8080`.

Running the app with `--help` will print out the available options:

```
./card-games-api --help
Card Games REST API

Usage:
  card-games-api [flags]

Flags:
      --bind-host string   Bind to hostname.
      --bind-port int      Listen on port. (default 8080)
      --db-type string     Database type: postgres or sqlite. (default "postgres")
      --db-url string      URL to sqlite database or PostgreSQL DSN.
  -h, --help               help for card-games-api
```

# Endpoints

## Deck Service
For the deck resource the following endpoints are available:
### CreateDeck

Creates a new deck of cards.
Based on the provided query parameters it will create either a full deck of cards or a partial one, shuffled or in order.

* Method: `POST`
* Path: `/v1/deck`
* Query Params:
  * `shuffled` - *optional*, boolean value. If set to `true`, the created deck will be shuffled.
  * `cards` - *optional*, list of cards as comma-separated string. If supplied will create a partial deck with the given cards. The cards must be valid and not duplicated.
  If not supplied, it will create a full deck of 52 cards.

**Examples**

Create a full deck in order:
```bash
export HOST=http://localhost:8080

curl -X POST "${HOST}/v1/deck"

{
  "deck_id": "ed7cfe37-ca0f-4216-884b-4a7442449c4b",
  "shuffled": false,
  "remaining": 52
}

```

Create a full shuffled deck:
```bash
export HOST=http://localhost:8080

curl -X POST "${HOST}/v1/deck?shuffled=true"

{
  "deck_id": "a0af8e56-023d-48ab-a7f6-7e9b10d56bae",
  "shuffled": true,
  "remaining": 52
}

```

Create a partial deck:
```bash
export HOST=http://localhost:8080

curl -X POST "${HOST}/v1/deck?cards=AC,2C,3C"

{
  "deck_id": "47eb9fb4-eadc-440b-9680-7be1ee225cf9",
  "shuffled": false,
  "remaining": 3
}
```

Try to create a partial deck with invalid card:
```bash
export HOST=http://localhost:8080

curl -X POST "${HOST}/v1/deck?cards=AC,2C,JJ,3C"
400
{
  "message": "invalid cards values: JJ"
}
```

### OpenDeck

Opens a deck - show all remaining cards in the deck.

* Method: `GET`
* Path: `/v1/deck/{deckId}`
* Parameter:
  * `deckId` - the ID of the deck to be opened.

**Examples**

Open a deck:
```bash
export HOST=http://localhost:8080
export DECK="47eb9fb4-eadc-440b-9680-7be1ee225cf9"  # Previously created deck ID.

curl "${HOST}/v1/deck/${DECK}"

{
  "deck_id": "47eb9fb4-eadc-440b-9680-7be1ee225cf9",
  "shuffled": false,
  "remaining": 3,
  "cards": [
    {
      "value": "ACE",
      "suit": "CLUBS",
      "code": "AC"
    },
    {
      "value": "2",
      "suit": "CLUBS",
      "code": "2C"
    },
    {
      "value": "3",
      "suit": "CLUBS",
      "code": "3C"
    }
  ]
}

```

Try to open a non-existing deck:
```bash
export HOST=http://localhost:8080
export DECK="47eb9fb4-eadc-440b-0000-0000000000000"  # Non existing deck ID.

curl "${HOST}/v1/deck/${DECK}"

404
{
  "message": "no such deck"
}
```

### DrawCards

Draws a number of cards from the deck.
Once the card is drawn, it is no longer in the deck.
Trying to draw more cards than there are in the deck will result in an error 400.

* Method: `POST`
* Path: `/v1/deck/{deckID}/draw`
* Path Parameter:
  * `deckId` - the ID of the deck to draw cards from
* Query Parameter:
  * `count` - *integer*, number of cards to draw from the deck

**Examples**

Draw one card from the deck:

```bash

export HOST=http://localhost:8080
export DECK="47eb9fb4-eadc-440b-9680-7be1ee225cf9"  # Previously created deck ID.

curl -X POST "${HOST}/v1/deck/${DECK}/draw?count=1"

{
  "cards": [
    {
      "value": "ACE",
      "suit": "CLUBS",
      "code": "AC"
    }
  ]
}

# Now display the deck
curl "${HOST}/v1/deck/${DECK}"

{
  "deck_id": "47eb9fb4-eadc-440b-9680-7be1ee225cf9",
  "shuffled": false,
  "remaining": 2,
  "cards": [
    {
      "value": "2",
      "suit": "CLUBS",
      "code": "2C"
    },
    {
      "value": "3",
      "suit": "CLUBS",
      "code": "3C"
    }
  ]
}

# The card is not in the deck anymore.
```

Draw multiple cards:

```bash

export HOST=http://localhost:8080
export DECK="ed7cfe37-ca0f-4216-884b-4a7442449c4b"  # The ID of the full deck that we created earlier.

curl -X POST "${HOST}/v1/deck/${DECK}/draw?count=4"

{
  "cards": [
    {
      "value": "ACE",
      "suit": "CLUBS",
      "code": "AC"
    },
    {
      "value": "2",
      "suit": "CLUBS",
      "code": "2C"
    },
    {
      "value": "3",
      "suit": "CLUBS",
      "code": "3C"
    },
    {
      "value": "4",
      "suit": "CLUBS",
      "code": "4C"
    }
  ]
}


# Now display the deck
curl "${HOST}/v1/deck/${DECK}"

{
  "deck_id": "ed7cfe37-ca0f-4216-884b-4a7442449c4b",
  "shuffled": false,
  "remaining": 48,
  "cards": [
    {
      "value": "5",
      "suit": "CLUBS",
      "code": "5C"
    },
    {
      "value": "6",
      "suit": "CLUBS",
      "code": "6C"
    },
    # ... omitted for brevity
    {
      "value": "QUEEN",
      "suit": "SPADES",
      "code": "QS"
    },
    {
      "value": "KING",
      "suit": "SPADES",
      "code": "KS"
    }
  ]
}

# The cards are not in the deck anymore.
```

Try to overdraw cards:

```bash

export HOST=http://localhost:8080
export DECK="47eb9fb4-eadc-440b-9680-7be1ee225cf9"  # The deck has only 2 cards left.

# Try to draw 5 cards.
curl -X POST "${HOST}/v1/deck/${DECK}/draw?count=5"

{
  "message": "not enough cards in deck"
}


# Now display the deck - no additional cards are drawn.
curl "${HOST}/v1/deck/${DECK}"

{
  "deck_id": "47eb9fb4-eadc-440b-9680-7be1ee225cf9",
  "shuffled": false,
  "remaining": 2,
  "cards": [
    {
      "value": "2",
      "suit": "CLUBS",
      "code": "2C"
    },
    {
      "value": "3",
      "suit": "CLUBS",
      "code": "3C"
    }
  ]
}

```