clean: 
	rm -rf generatede

generate: api.yml
	@echo "Generating files..."
	mkdir -p generated
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

init: clean generate
	go mod tidy

test:
	mockery --all --recursive --dir core --output core/mocks
	go clean -testcache
	go test -short -coverprofile=coverage.out -v ./core/usecase/...

test_api:
	go clean -testcache
	go test ./tests/...