DEPLOY_PATH := deploy/test

TASK_TRACKER_PATH := task-tracker
ACCOUNTING_PATH := accounting
ANALYTICS_PATH := analytics

install-all:
	make install -C auth
	make install -C task-tracker

run-infra:
	@echo '=====> Run containers...'
	@docker compose \
		-f ${DEPLOY_PATH}/infra.yml \
		-f ${DEPLOY_PATH}/storages.yml \
	up -d
	@echo
	@echo '=====> Apply migrations to task-tracker...'
	@make -C ${TASK_TRACKER_PATH} migrations-up
	@echo '=====> Apply migrations to accounting...'
	@make -C ${ACCOUNTING_PATH} migrations-up
	@#echo '=====> Apply migrations to analytics...'
	@#make -C ${ANALYTICS_PATH} migrations-up
	@echo

