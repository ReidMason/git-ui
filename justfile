run-test:
  go run . testfile.txt
test-cover-html:
  go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
