run:
	go run cmd/bagicore/main.go

dev:
	go mod download

generate-openapi:
	oapi-codegen -config=others/oapi-codegen/config.yaml \
		-templates others/oapi-codegen/templates \
		-package rapi \
		-o internal/controller/chttp/rapi/api.gen.go \
		api/openapi.yaml
