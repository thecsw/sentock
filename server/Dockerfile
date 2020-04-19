FROM golang
WORKDIR /app
COPY . .
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/lib/pq
RUN go get -u github.com/sanity-io/litter
RUN go build
CMD ./app