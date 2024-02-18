build-docker:
	docker build -t http-server .

build-docker-gosu:
	docker build -t http-server -f Dockerfile.gosu .

run-server:
	docker run --rm --name server-container -p 8085:8080 http-server
