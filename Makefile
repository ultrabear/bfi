GO=go build

bfi: main.go
	${GO} main.go

install: bfi
	cp ./bfi /bin/bfi
