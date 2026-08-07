.PHONY: buf-lint buf-build buf-generate

buf-lint:
	buf lint

buf-build:
	buf build

buf-generate:
	buf generate

push:
	buf push
