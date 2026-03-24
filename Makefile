run:
	go run cmd/mailer/*.go
build-linux:
	rm -f build/* && cd cmd/mailer && go build && cd - && mkdir -p build && mv cmd/mailer/mailer build/
build-windows:
	rm -f build/* && cd cmd/mailer && GOOS=windows GOARCH=amd64 go build && cd - && mkdir -p build && mv cmd/mailer/mailer.exe build/

