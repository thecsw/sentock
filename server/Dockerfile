FROM golang
WORKDIR /app
COPY . .
RUN go get -v -u ./...
RUN go build -v
CMD ./server