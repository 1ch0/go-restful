#include makefiles/const.mk
#include makefiles/dependency.mk
#include makefiles/release.mk
#include makefiles/develop.mk
#include makefiles/build.mk
#include makefiles/e2e.mk

.DEFAULT_GOAL := all
all: build

.PHONY: run-apiserver
run:
	go run ./cmd/apiserver/main.go

# Run tests
test: unit-test-core test-cli-gen
	@$(OK) unit-tests pass

unit-test-apiserver:
	go test -gcflags=all=-l -coverprofile=coverage.txt $(shell go list ./pkg/... ./cmd/...  | grep -E 'apiserver')

.PHONY: docker-build-apiserver
docker-build-apiserver:
	docker build --build-arg=VERSION=$(API_VERSION) --build-arg=GITVERSION=$(GIT_COMMIT) -t $(API_APISERVER_IMAGE) -f Dockerfile.apiserver .

build-cleanup:
	rm -rf _bin

# Run go fmt against code
fmt: goimports installcue
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

staticcheck: staticchecktool
	$(STATICCHECK) ./...

lint: golangci
	$(GOLANGCILINT) run ./...

reviewable:  fmt vet lint staticcheck
	go mod tidy

# Execute auto-gen code commands and ensure branch is clean.
check-diff: reviewable
	git --no-pager diff
	git diff --quiet || ($(ERR) please run 'make reviewable' to include all changes && false)
	@$(OK) branch is clean

# Push the docker image
docker-push:
	docker push $(API_CORE_IMAGE)

build-swagger:
	go run ./cmd/apiserver/main.go build-swagger ./docs/apidoc/swagger.json

image-cleanup:
ifneq (, $(shell which docker))
# Delete Docker images

ifneq ($(shell docker images -q $(API_CORE_TEST_IMAGE)),)
	docker rmi -f $(API_CORE_TEST_IMAGE)
endif

ifneq ($(shell docker images -q $(API_RUNTIME_ROLLOUT_TEST_IMAGE)),)
	docker rmi -f $(API_RUNTIME_ROLLOUT_TEST_IMAGE)
endif

endif


# Run tests
core-test:
	go test ./pkg/... -coverprofile cover.out


#BINARY = go-restful-apiserver
#IMAGE_REGISTRY =
#IMAGE_REPO =
#IMAGE_TAG = latest
#
#.PHONY: vendor build image release
#
#.PHONY: all
#all: check build
#.PHONY: run
#run:
#	go run ./cmd/apiserver/main.go
#
#.PHONY: check
#check:
#	go fmt ./...
#	go vet ./...
#.PHONY: build
#build:
#	@go build -o $(BINARY) ./cmd/apiserver/main.go
#.PHONY: clean
#clean:
#	rm -f $(BINARY)
#.PHONY: lint
#lint:
#	golangci-lint run -v --config ./.golangci.yml --timeout 5m
#.PHONY: test
#test:
#	go test ./...
#.PHONY: cover
#cover:
#	go test ./... -coverprofile coverage.out
#	go tool cover -html=coverage.out
#	rm -f coverage.out
#.PHONY: vendor
#vendor:
#	go mod tidy && go mod vendor
#
#build:
#	go build -o $(BINARY) cmd/server/main.go
#.PHONY: image
#image:
#	docker build -t $(IMAGE_REGISTRY)/$(IMAGE_REPO):$(IMAGE_TAG) .
#.PHONY: release
#release: image
#	docker push $(IMAGE_REGISTRY)/$(IMAGE_REPO):$(IMAGE_TAG)
#
### swagger: Generate swagger document.
#.PHONY: swagger
#swagger:
#	swag init -g ./cmd/server/main.go -d ./
#
#.PHONY: help
#help:
#	@echo "============================================="
#	@echo "make all     格式化go代码 并编译生成二进制文件"
#	@echo "make build   编译go代码生成二进制文件"
#	@echo "make clean   清理中间目标文件"
#	@echo "make test    执行测试case"
#	@echo "make check   格式化go代码"
#	@echo "make cover   检查测试覆盖率"
#	@echo "make run     直接运行程序"
#	@echo "make lint    执行代码检查"
#	@echo "make image   构建docker镜像"
#	@echo "make release 推送docker镜像"
#	@echo "make swagger 生成 swagger 接口文档"
#	@echo "============================================="

