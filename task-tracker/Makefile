ARTEFACT_NAME=task-tracker

GOOSE_DRIVER=postgres

ifndef GOOSE_DBSTRING
override GOOSE_DBSTRING = "postgres://postgres:password@localhost:5433/postgres?sslmode=disable"
endif

.PHONY: install
install:
	go mod tidy
	go build -o ${GOBIN}/${ARTEFACT_NAME} cmd/*.go
	chmod +x ${GOBIN}/${ARTEFACT_NAME}

migrations-status:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	goose -dir ./migrations status

migrations-up:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	goose -dir ./migrations up

migrations-down:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	goose -dir ./migrations down

generate-schemas:
	protoc --go_out=. --proto_path=../schemas/events ../schemas/events/task-workflow.proto