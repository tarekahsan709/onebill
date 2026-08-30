APP             := onebill
DOCKER_IMAGE    := ghcr.io/tarekahsan709/onebill
VERSION         := $(shell git describe --tags --exact-match 2>/dev/null)
LD_FLAGS        := "-s -w -X 'onebill.app/internal/version.Version=$(VERSION)'"
PKG_LIST        := $(shell go list ./... | grep -v /vendor/)
DB_URL          := postgres://postgres:postgres@localhost/onebill_test?sslmode=disable

export PGPASSWORD := postgres

.PHONY: \
	onebill \
	onebill-no-pie \
	linux-amd64 \
	linux-arm64 \
	linux-armv7 \
	linux-armv6 \
	linux-armv5 \
	linux-riscv64 \
	darwin-amd64 \
	darwin-arm64 \
	freebsd-amd64 \
	openbsd-amd64 \
	build \
	run \
	clean \
	add-string \
	test \
	lint \
	integration-test \
	clean-integration-test \
	docker-image \
	docker-image-distroless \
	docker-images

onebill:
	@ go build -buildmode=pie -ldflags=$(LD_FLAGS) -o $(APP)

onebill-no-pie:
	@ go build -ldflags=$(LD_FLAGS) -o $(APP)

linux-amd64:
	@ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

linux-arm64:
	@ CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

linux-armv7:
	@ CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

linux-armv6:
	@ CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

linux-armv5:
	@ CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

linux-riscv64:
	@ CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

darwin-amd64:
	@ GOOS=darwin GOARCH=amd64 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

darwin-arm64:
	@ GOOS=darwin GOARCH=arm64 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

freebsd-amd64:
	@ CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

openbsd-amd64:
	@ GOOS=openbsd GOARCH=amd64 go build -ldflags=$(LD_FLAGS) -o $(APP)-$@
	@ sha256sum $(APP)-$@ > $(APP)-$@.sha256

build: linux-amd64 linux-arm64 linux-armv7 linux-armv6 linux-armv5 linux-riscv64 darwin-amd64 darwin-arm64 freebsd-amd64 openbsd-amd64

run:
	@ LOG_DATE_TIME=1 LOG_LEVEL=debug RUN_MIGRATIONS=1 CREATE_ADMIN=1 ADMIN_USERNAME=admin ADMIN_PASSWORD=test123 go run main.go

clean:
	@ rm -f $(APP)-* $(APP) $(APP)*.exe $(APP)*.sha256

add-string:
	cd internal/locale/translations && \
	for file in *.json; do \
		jq --indent 4 --arg key "$(KEY)" --arg val "$(VAL)" \
		   '. + {($$key): $$val} | to_entries | sort_by(.key) | from_entries' "$$file" > tmp && \
		mv tmp "$$file"; \
	done

test:
	go test -cover -race -count=1 ./...

lint:
	go vet ./...
	test -z "$$(gofmt -l .)"
	golangci-lint run

integration-test:
	psql -U postgres -c 'drop database if exists onebill_test;'
	psql -U postgres -c 'create database onebill_test;'

	DATABASE_URL=$(DB_URL) \
	ADMIN_USERNAME=admin \
	ADMIN_PASSWORD=test123 \
	CREATE_ADMIN=1 \
	RUN_MIGRATIONS=1 \
	LOG_LEVEL=debug \
	go run main.go >/tmp/onebill.log 2>&1 & echo "$$!" > "/tmp/onebill.pid"

	while ! nc -z localhost 8080; do sleep 1; done

	TEST_ONEBILL_BASE_URL=http://127.0.0.1:8080 \
	TEST_ONEBILL_ADMIN_USERNAME=admin \
	TEST_ONEBILL_ADMIN_PASSWORD=test123 \
	go test -v -count=1 ./internal/api

clean-integration-test:
	@ kill -9 `cat /tmp/onebill.pid`
	@ rm -f /tmp/onebill.pid /tmp/onebill.log
	@ psql -U postgres -c 'drop database if exists onebill_test;'

docker-image:
	docker build --pull -t $(DOCKER_IMAGE):$(VERSION) -f packaging/docker/alpine/Dockerfile .

docker-image-distroless:
	docker build -t $(DOCKER_IMAGE):$(VERSION) -f packaging/docker/distroless/Dockerfile .

docker-images:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--file packaging/docker/alpine/Dockerfile \
		--tag $(DOCKER_IMAGE):$(VERSION) \
		--push .
