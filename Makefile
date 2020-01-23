all: restore-deps test vet

test:
	go test -v ./minidom
	
vet: 
	go vet ./minidom

clean:
	rm *.test

restore-deps:
	go mod tidy