FROM golang
RUN apt-get update
RUN apt-get install -y postgresql-client

RUN mkdir /go/src/goiex
WORKDIR /go/src/goiex
ENV GOBIN=/go/bin
COPY src/go* .
RUN go install .

CMD /go/bin/goiex
