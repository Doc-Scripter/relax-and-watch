run:
	clear
	cd backend/cmd && go run main.go
format:
	cd backend/cmd && gofmt -w -s .
restart-server:
	fuser -k 8080/tcp&&run