build:
	docker build -t bolgsrv:1 .
run:
	docker run -itd -p 8080:8080 --name bolgsrv bolgsrv:1
stop:
	docker rm -f bolgsrv
restart:
	make stop
	make run
all:
	make build
	make restart
logs:
	docker logs -f bolgsrv
