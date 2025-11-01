.PHONY: build clean docker start stop

build:
	go build -o bin/dive-beacon .

clean:
	rm -rf bin/

docker:
	docker build -t dive-beacon .

start:
	docker compose up -d

stop:
	docker compose down