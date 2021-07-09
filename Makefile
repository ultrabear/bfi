GOCMD=go
EXE=bfi

${EXE}: *.go */*.go
	${GOCMD} build -o ${EXE}

install: ${EXE}
	cp ./${EXE} /bin/${EXE}

userinstall: ${EXE}
	cp ./${EXE} ~/.local/bin/${EXE}
