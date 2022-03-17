build:
	docker build -t bolgsrv:v1 .
run:
	docker run -itd -p 8080:8080 --name bolgsrv bolgsrv:v1
stop:
	docker rm -f bolgsrv
restart:
	make stop
	make run
logs:
	docker logs -f bolgsrv
