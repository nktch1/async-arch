include auth/Makefile

install-all:
	make install -C auth
	make install -C task-tracker

run-in-docker:
	docker compose \
		-f deploy/test/docker-compose.yml \
		up -d --force-recreate --build
