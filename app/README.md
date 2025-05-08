The `tr-extractor` backend is written in Golang. It is a combination of API Endpoints and background processing.

## Go Module

```bash
go mod init github.com/khaledhikmat/tr-extractor
go get -u github.com/joho/godotenv
go get -u github.com/jmoiron/sqlx
go get -u github.com/lib/pq
go get -u github.com/gin-gonic/gin 
go get -u github.com/gin-contrib/cors
go get -u github.com/mdobak/go-xerrors
go get -u github.com/fatih/color
go get -u go.opentelemetry.io/otel
go get -u go.opentelemetry.io/contrib/exporters/autoexport
go get -u go.opentelemetry.io/contrib/propagators/autoprop
go get -u github.com/aws/aws-sdk-go
go get -u github.com/aws/aws-sdk-go-v2
go get -u github.com/aws/aws-sdk-go-v2/config
go get -u github.com/aws/aws-sdk-go-v2/service/s3
go get -u github.com/aws/aws-sdk-go-v2/feature/s3/manager
go get -u github.com/joho/godotenv
go get -u github.com/google/uuid
```

## Env Variables

| NAME           | DEFAULT | DESCRIPTION       |
|----------------|-----|------------------|
| TRELLO_API_KEY       | `trello-api-key`  | Trello API Key. |
| TRELLO_TOKEN       | `trello-token`  | Trello Token. |
| TRELLO_SECRET       | `trello-secret`  | Trello Secret. |
| TRELLO_BASE_URL       | `trello-base-url`  | Trello Base URL. |
| TRELLO_PROPERTIES_BOARD_ID       | `trello-properties-board-id`  | Trello Properties Board ID. |
| TRELLO_EXPENSES_BOARD_ID       | `trello-expenses-board-id`  | Trello Expenses Board ID. |
| TRELLO_INHERITANCE_CONFINEMENTS_BOARD_ID       | `trello-inheritance-confinements-board-id`  | Trello Inheritance and Confinements Board ID. |
| TRELLO_TODO_BOARD_ID       | `trello-todo-board-id`  | Trello TODO Board ID. |
| DB_DSN       | `railway-postgres-db`  | HTTP Server port. Required to expose API Endpoints. |
| PROPERTIES_EXCEL_UPDATE_WEBHOOK       | `empty`  | Webhook URL for update Google properties sheet |
| PROPERTIES_NOTION_UPDATE_WEBHOOK       | `empty`  | Webhook URL for update Notion properties database |
| APP_NAME       | `tr-extractor`  | Name of the microservice to appear in OTEL. |
| API_PORT       | `8080`  | HTTP Server port. Required to expose API Endpoints. |
| RUN_TIME_ENV  | `dev`  | Runetime env name.  |
| OPEN_TELEMETRY     | `false`  | If `true`, it disables collecting OTEL telemetry signals.   |
| OTEL_EXPORTER_OTLP_ENDPOINT     | `http://localhost:4318`  | OTEL endpoint.   |
| OTEL_SERVICE_NAME     | `yt-extractor-backend`  | OTEL application name.   |
| OTEL_GO_X_EXEMPLAR     | `true`  | OTEL GO.   |

## Run Locally

```bash
go run main.go
```

## Build and Push to Docker Hub

```bash
make push-2-hub
```

## Docker Prune

```bash
# Remove all stopped containers
docker container prune -f

# Remove all unused images
docker image prune -a -f

# Remove all unused volumes
docker volume prune -f

# Remove all unused networks
docker network prune -f

# Remove all unused data
docker system prune -a -f --volumes
```


