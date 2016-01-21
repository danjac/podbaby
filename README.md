Podbaby is a web application for listing to and organizing your podcasts.

Development setup
=================

We assume the following are installed:

- Go 1.5+
- PostgreSQL 8+
- NodeJS 5.4+

Install
=======

```
go get github.com/danjac/podbaby
cd GOPATH/src/github.com/danjac/podbaby
cp .env.sample .env // edit as required
make
make test
./bin/runapp serve -env=dev 
npm run dev
```

Open browser at localhost:5000


