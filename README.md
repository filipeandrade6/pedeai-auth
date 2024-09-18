GOARCH=amd64 GOOS=linux go build -o bootstrap main.go
zip deployment.zip bootstrap