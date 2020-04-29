FROM golang
WORKDIR /app
COPY . .
RUN go get golang.org/x/tools/cmd/godoc
RUN godoc -http=:6060 &
#RUN wget -m -k -q -erobots=off --no-host-directories --no-use-server-timestamps http://localhost:6060
RUN go get -v ./...
RUN go build -v
