ARTEFACT_NAME=task-tracker

.PHONY: install
install:
	go mod tidy
	go build -o ${GOBIN}/${ARTEFACT_NAME} cmd/*.go
	chmod +x ${GOBIN}/${ARTEFACT_NAME}
