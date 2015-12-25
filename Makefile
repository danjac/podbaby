GOPATH := ${PWD}:${GOPATH}
export GOPATH

build: build-server build-ui

build-server: 
	godep restore
	go build -o bin/serve -i main.go

build-ui:
	npm install
	npm run build

	
clean-ui:
	rm -rf node_modules

test:
	go test ./...
	npm run test

clean: 
	clean-ui
