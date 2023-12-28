.PHONY: gen-sqlc gen-proto

gen-proto:
	docker compose up generate_pb_go
gen-sqlc:
	docker compose up generate_sqlc