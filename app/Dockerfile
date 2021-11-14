FROM golang
RUN apt-get update
RUN apt-get install -y postgresql-client

RUN mkdir /go/src/goiex
WORKDIR /go/src/goiex
COPY src/go.mod .
COPY src/go.sum .
COPY bin/goiex /go/bin/
RUN go get github.com/lib/pq

CMD tail -f /dev/null
