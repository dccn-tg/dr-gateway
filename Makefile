ifndef GOPATH
	GOPATH := $(HOME)/go
endif

ifndef GOOS
	GOOS := linux
endif

ifndef GO111MODULE
	GO111MODULE := on
endif

all: build

build: api-server

api-server:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -o bin/api-server internal/api-server/main.go

swagger:
	swagger validate pkg/swagger/swagger.yaml
	go generate github.com/dccn-tg/dr-gateway/internal/api-server github.com/dccn-tg/dr-gateway/pkg/swagger

doc: swagger
	swagger serve pkg/swagger/swagger.yaml

test_dr: build
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) DR_GATEWAY_CONFIG=$(DR_GATEWAY_CONFIG) go test -v github.com/dccn-tg/dr-gateway/pkg/dr/... 

clean:
	rm -rf bin
