# Example Application

This is an example application that demostrates receiving protobuf bytes and returning json.

## generate proto

`protoc -I=./pb --go_out=./pb --go_opt=paths=source_relative ./pb/*.proto`

## Run application

`go run main.go`
