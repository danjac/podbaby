GOPATH := ${PWD}:${GOPATH}
export GOPATH

build: build-service build-ui

build-service: 
	godep restore
	go build -o bin/serve -i main.go

build-ui:
	npm install
	npm run build
	
clean-ui:
	rm -rf node_modules

clean: 
	clean-ui

test-service:
	go test ./...

test-ui:
	npm run test

test: test-service test-ui

