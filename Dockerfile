FROM golang
WORKDIR /server
COPY . .
RUN go get .
RUN protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/upload_file.proto
CMD go run main.go