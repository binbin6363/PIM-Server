.PHONY: build build-all clean vet fmt run

GOCMD=GO111MODULE=on go
BINARY=bin/im_server

build:
	@echo "build go binary..."
	${GOCMD} build -o ${BINARY} server/main.go
	@echo "build go binary done"

build-all:
	@echo "build all go files..."
	${GOCMD} build ./...
	@echo "build all go files done"

clean:
	@echo "clean go binary..."
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@echo "clean go binary done"

vet:
	@echo "vet go files..."
	${GOCMD} vet ./...
	@echo "vet go files done"

fmt:
	@echo "fmt go files..."
	${GOCMD} fmt ./...
	@echo "fmt go files done"

run:
	@echo "run server..."
	${GOCMD} run server/main.go -f etc/conf.yaml || true
