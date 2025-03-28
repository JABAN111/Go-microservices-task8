container_runtime := $(shell which docker || which podman)

$(info using ${container_runtime})

up: down
	${container_runtime} compose up --build -d

down:
	${container_runtime} compose down

clean:
	${container_runtime} compose down -v

run-tests: 
	${container_runtime} run --rm --network=host tests:latest

test:
	make clean
	make up
	@echo wait cluster to start && sleep 10
	make run-tests
	make clean
	@echo "test finished"

lint:
	make -C search-service lint

proto:
	make -C search-service protobuf
