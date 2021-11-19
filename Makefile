# Varibles
GOIEX_HOME=`pwd`
IEX_APP=app/bin/goiex
APP_DOCKER_FILE=app/Dockerfile
PG_DOCKER_FILE=pg/Dockerfile
IMAGES=goiex-postgres-centos goiex-app
# Rules
all: build down up

build: build_app build_docker

build_app: app/src/Makefile
	make -C app/src

build_docker: ${IMAGES}

goiex-postgres-centos: ${PG_DOCKER_FILE}
	(cd pg; docker build --rm --progress=plain --build-arg DB_USER=${DB_USER} --build-arg DB_PASSWORD=${DB_PASSWORD} --build-arg DB_NAME=${DB_NAME} -t $@ .)

goiex-app: ${APP_DOCKER_FILE}
	(cd app; docker build --rm --progress=plain -t $@ .)

up: docker-compose.yml .env
	docker-compose up -d

down: docker-compose.yml .env
	docker-compose down

clean:
	rm -rf ${IEX_APP}

login:
	docker exec -it goiex_goiex-app_1 bash