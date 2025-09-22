test:
	go test -v -cover

cover:
	go test -v -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

profiler:
	go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html: