include ./env/local.env

comma := ,
space := $(empty) $(empty)


docker-build:
	docker compose build

docker-start:
	docker compose up

docker-clean:
	docker compose down --volumes --remove-orphans

build-all: 
	make clean
	@for service in $(subst $(comma),$(space),$(SERVICES)); do \
		go build -o bin/$$service lexyblazy.github.com/microservices-starter/cmd/$$service; \
	done 

clean:
	rm -rf ./bin

start-services:
	@for service in $(subst $(comma),$(space),$(SERVICES)); do \
		echo "Starting $$service service ..."; \
		cat local.env | ./bin/$$service & \
	done

start: 
	echo "Starting $(service) service ..."
	cat local.env | ./bin/$$service

