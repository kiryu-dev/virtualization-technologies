build-docker:
	docker build -t http-server .

run-server:
	docker run --rm --name server-container -p 8085:8080 http-server
