.PHONY: all
all: ;

.PHONY: pg
pg:
	docker run --rm \
		--name=praktikum-webinar-db \
		-v $(abspath ./db/init/):/docker-entrypoint-initdb.d \
		-v $(abspath ./db/data/):/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD="P@ssw0rd" \
		-d \
		-p 5432:5432 \
		postgres:16.3

.PHONY: stop-pg
stop-pg:
	docker stop praktikum-webinar-db

.PHONY: clean-data
clean-data:
	sudo rm -rf ./db/data/

.PHONY: build-datagen
build-datagen:
	go build -o ./app/bin/datagen app/cmd/datagen

.PHONY: build-app
build-app:
	go build -o ./app/bin/employees app/cmd/employees
