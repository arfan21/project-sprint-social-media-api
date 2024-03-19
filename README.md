## Getting Started <a name = "getting_started"></a>

### Prerequisites

-   Setup env variable
-   Install dependencies

```bash
# swag command
# refer https://github.com/swaggo/swag
go install github.com/swaggo/swag/cmd/swag@latest

# mockery
# refer https://vektra.github.io/mockery/latest/installation/
go install github.com/vektra/mockery/v2@v2.40.1

# air (golang hot reload)
# refer https://github.com/cosmtrek/air
go install github.com/cosmtrek/air@latest
```

## Build <a name="build"></a>

```
go build -o server ./cmd/main.go
```

## Usage <a name="usage"></a>

### Run Server

```
./server serve
```

### Run Migration

```
./server migrate up
./server migrate down
./server migrate fresh
./server migrate drop
```

## Development <a name="development"></a>

### Create Migration

install [golang migrate](https://github.com/golang-migrate/migrate)

```
migrate create -ext sql <migration_name>
```

### Live Reload Server

install [air](https://github.com/cosmtrek/air)

```
air serve
```

for windows

```
make air-win
```

### Generate Mock

install [mockery](https://github.com/vektra/mockery)

```
make mocks
```

### Generate Swagger Docs

install [swaggo](https://github.com/swaggo/swag)

```
make swag
```

### Run Test Locally

```
go test ./... -v
```
