prepare-swagger-gen:
	docker pull redocly/cli

generate-swagger: prepare-swagger-gen
	docker run --rm -v $(PWD):/spec redocly/cli bundle api/schema.yaml -o api/bundle.yaml

generate-oapi:
	go tool oapi-codegen -config oapi.yaml api/bundle.yaml

generate-oapi-client:
	go tool oapi-codegen -config oapi-client.yaml api/bundle.yaml

generate-go:
	go generate ./...

generate-sqlc:
	go tool sqlc generate -f internal/repository/rooms/gen/sqlc.yaml
	go tool sqlc generate -f internal/repository/bookings/gen/sqlc.yaml

generate: generate-go generate-swagger generate-oapi generate-oapi-client generate-sqlc

test:
	go test -v ./...

test-short:
	go test --short -v ./...