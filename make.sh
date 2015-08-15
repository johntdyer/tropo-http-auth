GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pkg/tropo-http-auth tropo-http-auth.go
CGO_ENABLED=0 go build -o pkg/tropo-http-auth-osx tropo-http-auth.go
