FROM golang

ENV HOMEDIR /go/src/github.com/danjac/podbaby 
ENV DB_URL postgres://postgres@db/postgres?sslmode=disable

ADD . ${HOMEDIR}
WORKDIR ${HOMEDIR}

RUN echo $(cat /etc/hosts)
RUN curl -sL https://deb.nodesource.com/setup_0.10 | bash -
RUN apt-get install -y build-essential 
RUN apt-get install -y nodejs 
RUN go get github.com/mattes/migrate
RUN go get github.com/tools/godep
RUN make

ENTRYPOINT ["./docker-entrypoint.sh"]

