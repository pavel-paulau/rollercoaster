build:
	go build -v

fmt:
	gofmt -w -s *.go

test:
	go test -v -cover -race

docker:
	CGO_ENABLED=0 go build -v -a --ldflags "-s" && upx -q6 rollercoaster
	docker build --rm -t perflab/rollercoaster .

clean:
	rm -fr rollercoaster build
