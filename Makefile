PYROSCOPE_VERSION ?= latest
SERVICE_NAME      ?= pyroscope-smoke-test
SEARCH_TERM       ?= smokeTestMagicFunction
TIMEOUT           ?= 60

.PHONY: test
test: up
	docker run -d --name smoke-app \
		--network pyroscope-smoke_default \
		-e PYROSCOPE_SERVER_ADDRESS=http://pyroscope:4040 \
		-e PYROSCOPE_SERVICE_NAME=$(SERVICE_NAME) \
		smoke-app:local

	./querier/querier \
		--service-name "$(SERVICE_NAME)" \
		--term "$(SEARCH_TERM)" \
		--timeout "$(TIMEOUT)s"

.PHONY: up
up: build
	docker compose -p pyroscope-smoke -f internal/docker-compose.yml up -d
	@for i in $$(seq 1 30); do \
		curl -sf http://localhost:4040/ready && break; \
		echo "Waiting for Pyroscope... (attempt $$i)"; \
		sleep 2; \
	done

.PHONY: down
down:
	docker stop smoke-app
	docker rm smoke-app
	docker compose -p pyroscope-smoke -f internal/docker-compose.yml down --volumes

.PHONY: build
build:
	docker build -f testdata/app/Dockerfile . -t smoke-app
	cd querier && go build -o querier .

.PHONY: clean
clean:
	rm -f querier/querier
