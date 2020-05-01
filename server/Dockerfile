FROM golang
WORKDIR /app
COPY . .
RUN go get -v ./...
RUN go get -v golang.org/x/tools/cmd/godoc
RUN go build -v
