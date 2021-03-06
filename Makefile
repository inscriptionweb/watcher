path = $(PWD)/../../../../
src_path = $(PWD)

build:	clean
	@echo "*** building ***"
	@export GOPATH=$(path) && go build
clean:
	@echo "*** cleaning ***"
	@rm -rf watcher
	@gofmt -s -w .
test:
	@echo "*** tests ***"
	@export GOPATH=$(path) && cd $(src_path)/sender && go test
	@export GOPATH=$(path) && cd $(src_path)/tree_walker && go test
