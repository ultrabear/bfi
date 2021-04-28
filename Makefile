bfi: main.go
	go build

install: bfi
	cp ./bfi /bin/bfi
