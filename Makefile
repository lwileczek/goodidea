build:
	mkdir -p static/img
	npm run build
	go generate ./...
	go build -o server ./app/main.go 

run:
	mkdir -p static/img
	npm run build
	go generate ./...
	go run ./app/main.go

deployLambda:
	bash utils/deploy.sh

buildContainer:
	npm run build
	go generate ./...
	mkdir app/kodata
	cp -r static templates app/kodata
	ko build ./app/...
