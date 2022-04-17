.PHONY: .
-include .env
export
network:
	@docker network create ${TRAEFIK_NETWORK}
init:
	docker-compose up -d --force-recreate green db traefik
build:
	docker-compose build green blue
green:
	docker-compose up -d --force-recreate green
	docker-compose stop -t $(DOCKER_STOP_TIMOUT) blue
	docker-compose rm -f blue
blue:
	docker-compose up -d --force-recreate blue
	docker-compose stop -t $(DOCKER_STOP_TIMOUT) green
	docker-compose rm -f green