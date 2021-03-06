# Varibles
GOIEX_HOME=`pwd`
IEX_APP=app/bin/goiex
APP_DOCKER_FILE=app/Dockerfile
APP_DOCKER_FILE_DEV=app/Dockerfile.dev
PG13_DOCKER_FILE=pg13/Dockerfile
PG12_DOCKER_FILE=pg12/Dockerfile
NGINX_DOCKER_FILE=front/Dockerfile
NGINX_DOCKER_FILE_DEV=front/Dockerfile.dev

IMAGES=goiex-postgres12-centos goiex-app goiex-app-dev goiex-ngnix goiex-ngnix-dev
# Rules
all: build down up

build: build_app build_docker

build_app: app/src/Makefile
	make -C app/src

build_docker: ${IMAGES}

goiex-postgres13-centos: ${PG13_DOCKER_FILE}
	(cd pg13; docker build --rm --progress=plain --build-arg DB_USER=${DB_USER} --build-arg DB_PASSWORD=${DB_PASSWORD} --build-arg DB_NAME=${DB_NAME} -t $@ .)

goiex-postgres12-centos: ${PG12_DOCKER_FILE}
	(cd pg12; docker build --rm --progress=plain --build-arg DB_USER=${DB_USER} --build-arg DB_PASSWORD=${DB_PASSWORD} --build-arg DB_NAME=${DB_NAME} -t $@ .)

goiex-app: ${APP_DOCKER_FILE}
	(cd app; docker build --rm --progress=plain -t $@ .)

goiex-app-dev: ${APP_DOCKER_FILE_DEV}
	(cd app; docker build --rm --progress=plain -t $@ -f Dockerfile.dev .)

goiex-ngnix: ${NGINX_DOCKER_FILE}
	(cd front; docker build --rm --progress=plain -t $@ .)

goiex-ngnix-dev: ${NGINX_DOCKER_FILE_DEV}
	(cd front; docker build --rm --progress=plain -t $@ -f Dockerfile.dev .)

up: docker-compose.yml setenv
	docker-compose up -d

down: docker-compose.yml setenv
	docker-compose down

clean:
	rm -rf ${IEX_APP}

login_app:
	docker exec -it goiex_goiex-app_1 bash

login_db:
	docker exec -it goiex_goiex-postgres-server_1 bash

log_app:
	docker logs -f goiex_goiex-app_1

log_front:
	docker logs -f goiex_goiex-front_1