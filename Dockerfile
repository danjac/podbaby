FROM golang

ADD . /go/src/github.com/danjac/podbaby 

WORKDIR /go/src/github.com/danjac/podbaby 

RUN echo $(cat /etc/hosts)
RUN curl -sL https://deb.nodesource.com/setup_0.10 | bash -
RUN apt-get install -y build-essential 
RUN apt-get install -y nodejs 
RUN go get github.com/tools/godep
RUN make

# set environment keys
#CMD migrate -path=./migrations -db= up

CMD ./bin/runapp serve

