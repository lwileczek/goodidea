run:
	npm run build
	go run ./app/main.go

deployLambda:
	bash utils/deploy.sh
