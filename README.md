# Mapon Backend Test


This repository contains a solution to the Backend Test from [Mapon](https://mapon.com/).

Please modify `config.yml` with proper data. Or set sensible configuration into `.env` file.

## Usage

To execute:

```
make build
make run
```

To clean:

```
make clean
```

To interact with dev environment using docker-compose:
```
make develenv-up   # docker-compose up whole environment
make develenv-down # docker-compose down whole environment
make develenv-sh   # logs bash session inside develenv container
```

To run test:
```
make develenv-up
make develenv-sh
make test-acceptance # (once inside develenv)
```

# Configuration

The following code shows the yaml configuration needed by the application:

```yaml
server:
  host:
  port: 8080
log:
  level: DEBUG
mongo:
  host: localhost:27017
  user:
  password:
  database: mapon
  collections:
    users: users
session:
  cookie: mapon-session
  expires: 3600
mapon:
  url: https://mapon.com/api/v1
  key:
  endpoints:
    unit: unit/list.json
    route: route/list.json
```

This configuration can also be achieved by environment variables:

| **Environment Variable** | **Description**                             |
| :----------------------- | :------------------------------------------ |
| `SERVER_HOST`            | Server host                                 |
| `SERVER_PORT`            | Server port                                 |
| `LOG_LEVEL`              | Server log level                            |
| `MONGO_HOST`             | MongoDB host                                |
| `MONGO_USER`             | MongoDB user *(sensible data)*              |
| `MONGO_PASSWORD`         | MongoDB password *(sensible data)*          |
| `MONGO_DATABASE`         | MongoDB Database                            |
| `MONGO_COLLECTION_USERS` | MongoDB users' collection                   |
| `SESSION_COOKIE`         | Session cookie name                         |
| `SESSION_EXPIRES`        | Session cookie expiration time (in seconds) |
| `MAPON_URL`              | Mapon API base URL                          |
| `MAPON_KEY`              | Mapon Key *(sensible data)*                 |
| `MAPON_ENDPOINTS_UNIT`   | Mapon Units endpoint                        |
| `MAPON_ENDPOINTS_ROUTE`  | Mapon Units Route                           |

## Authors

- Ismael Taboada Rodero: [@ismtabo](https://github.com/ismtabo)
