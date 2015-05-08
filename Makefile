path = $(PWD)/../../../../

build:	clean
	@echo "*** building ***"
	@export GOPATH=$(path) && go build
clean:
	@echo "*** cleaning ***"
	@rm -rf watcher
	@find -name *.go -exec go fmt {} \;
