test:
	go test -v ./...

benchmark:
	go test -bench . -benchmem -count 5 -run none > bench/v$(V)

cov:
	go test -coverprofile=cover.out
	go tool cover -func=cover.out > coverage/v$(V)
	rm cover.out