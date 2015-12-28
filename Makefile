GOPATH := ${PWD}:${GOPATH}
export GOPATH

build: build-backend build-frontend

build-backend:
	cd serve; godep restore
	cd serve; go build -o ../bin/serve -i main.go

build-frontend:
	npm install
	npm run build

clean-frontend:
	rm -rf node_modules

clean-backend:
	rm -rf bin

clean:
	clean-frontend
	clean-backend

test-backend:
	go test ./...

test-frontend:
	npm run test

test: test-backend test-frontend
