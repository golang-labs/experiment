dynamic_delete_collection:
	env GOOS=linux GOARCH=amd64 go build -o bin/dynamic_delete_collection ./cmd/dynamic_delete_collection.go