.PHONY: generate generate-sqlc generate-goverter generate-easyjson

generate:
	$(MAKE) generate-sqlc
	$(MAKE) generate-goverter
	$(MAKE) generate-easyjson

generate-sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate

generate-goverter:
	bash scripts/generate-goverter.sh

generate-easyjson:
	bash scripts/generate-easyjson.sh
