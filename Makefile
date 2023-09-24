test:
	go test

test-with-coverage:
	if [ -e cover.out ]; then rm cover.out; fi
	go test -cover -coverprofile cover.out && go tool cover -html=cover.out