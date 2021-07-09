GOCMD=go
EXE=bfi

bfi: *.go */*.go
	${GOCMD} build -o ${EXE}

install: bfi
	cp ./${EXE} /bin/${EXE}

userinstall: bfi
	cp ./${EXE} ~/.local/bin/${EXE}
