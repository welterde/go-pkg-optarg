
all:
	make -C optarg install

test:
	make -C optarg test

clean:
	make -C optarg clean

format:
	gofmt -w .

