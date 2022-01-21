help: ## show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%s\033[0m|%s\n", $$1, $$2}' \
	| column -t -s '|'

build: fmt ## build the exe
	GOOS=windows go build

fmt: ## format code
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

test:
	richgo test ./...

build-manifest: # rebuild rsrc.syso
	# force linux tooling for rsrc (container defaults to GOOS=windows)
	# GOOS=linux go get github.com/akavel/rsrc
	GOOS=linux go install github.com/akavel/rsrc
	rsrc -manifest pick-a-browser.exe.manifest -o rsrc.syso
