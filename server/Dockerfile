FROM golang
WORKDIR /app
COPY . .
RUN go get -v ./...
RUN go build -v
CMD ./server