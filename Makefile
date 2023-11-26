build:
	mkdir -p static/img
	npm run build
	go build -o server ./app/main.go 

run:
	mkdir -p static/img
	npm run build
	go run ./app/main.go

deployLambda:
	bash utils/deploy.sh

buildContainer:
	npm run build
	mkdir app/kodata
	cp -r static templates app/kodata
	ko build ./app/...
