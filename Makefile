install-all:
	make install -C auth
	make install -C task-tracker

run-in-docker:
	docker compose \
		-f deploy/test/docker-compose.yml \
		up -d --force-recreate --build

run-infra:
	docker compose \
	-f deploy/test/docker-compose.yml \
	up -d  zookeeper broker ates-schema-registry ates-accounting-db ates-task-tracker-db
