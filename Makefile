build_all:
	go build -v ./single_thread_blocking_io
	go build -v ./multi_thread_blocking_io
	go build -v ./non_blocking_busy_waiting
	go build -v ./single_thread_event_loop

test_all:
	go test -v ./single_thread_blocking_io/...
	go test -v ./multi_thread_blocking_io/...
	go test -v ./non_blocking_busy_waiting/...
	go test -v ./single_thread_event_loop/...

clean_all:
	cd single_thread_blocking_io && go clean -testcache && cd ..
	cd multi_thread_blocking_io && go clean -testcache && cd ..
	cd non_blocking_busy_waiting && go clean -testcache && cd ..
	cd single_thread_event_loop && go clean -testcache && cd ..

clean_test_all: clean_all test_all