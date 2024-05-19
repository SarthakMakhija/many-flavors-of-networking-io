build_all:
	go build -v ./single_threaded_blocking_io
	go build -v ./multi_threaded_blocking_io
	go build -v ./non_blocking_busy_waiting

test_all:
	go test -v ./single_threaded_blocking_io/...
	go test -v ./multi_threaded_blocking_io/...
	go test -v ./non_blocking_busy_waiting/...

clean_all:
	cd single_threaded_blocking_io && go clean -testcache && cd ..
	cd multi_threaded_blocking_io && go clean -testcache && cd ..
	cd non_blocking_busy_waiting && go clean -testcache && cd ..

clean_test_all: clean_all test_all