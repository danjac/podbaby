FROM golang

ADD . /go/src/github.com/danjac/podbaby 

WORKDIR /go/src/github.com/danjac/podbaby 

RUN echo $(cat /etc/hosts)
RUN curl -sL https://deb.nodesource.com/setup_0.10 | bash -
RUN apt-get install -y build-essential 
RUN apt-get install -y nodejs 
RUN go get github.com/tools/godep
RUN make

RUN go get github.com/mattes/migrate
RUN migrate -path=./migrations -url=postgres://postgres@db/postgres?sslmode=disable up

CMD ./bin/runapp serve 
