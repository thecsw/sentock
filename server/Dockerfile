FROM golang
WORKDIR /app
COPY . .
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/lib/pq
RUN go build
CMD ./app