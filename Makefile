GOCMD=go
EXE=odb-obsauth

${EXE}: *.go */*.go
	${GOCMD} build -o ${EXE}

tidy:
	gofmt -s -w .
	go mod tidy
	go vet

clean:
	rm ${EXE}
