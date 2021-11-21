FROM golang
RUN apt-get update
RUN apt-get install -y postgresql-client
#RUN pip3 install libpq-dev python3-dev gettext Pyrseas

RUN mkdir /go/src/goiex
WORKDIR /go/src/goiex
ENV GOBIN=/go/bin
ADD src .

RUN go install .
CMD /go/bin/goiex
