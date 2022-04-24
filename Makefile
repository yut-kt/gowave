test:
	go test -v ./...

benchmark:
	go test -bench . -benchmem -count 5 -run none > docs/bench/v$(V)

cov:
	go test -coverprofile=cover.out
	go tool cover -func=cover.out > docs/coverage/v$(V)
	rm cover.out

generate:
	jupyter nbconvert --to notebook \
 		--inplace \
 		--ClearMetadataPreprocessor.enabled=True \
 		--ClearMetadataPreprocessor.clear_notebook_metadata=False \
	 	--execute gen/sample_json_generator.ipynb
	go run gen/samples_generator.go