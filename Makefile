GOPATH := ${PWD}:${GOPATH}
export GOPATH

build: build-backend build-frontend

build-backend:
	godep restore
	go build -o ./bin/runapp -i main.go

build-frontend:
	npm install
	npm run build

clean-frontend:
	rm -rf node_modules

clean-backend:
	rm -rf bin

clean: clean-frontend clean-backend

test-backend:
	go test ./...

test-frontend:
	npm run test

test: test-backend test-frontend
