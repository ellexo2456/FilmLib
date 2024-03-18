go test -coverprofile=coverage.out.tmp ./...
grep -v "mocks" coverage.out.tmp | grep -v "docs" | grep -v "internal/app" | grep -v "swagger_types" > coverage.out
go tool cover -func=coverage.out
go tool cover --html=coverage.out
