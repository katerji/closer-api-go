FROM golang
WORKDIR /server
COPY . .
RUN go get .
CMD go run main.go