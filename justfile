default:
  just --list
 

lint: build ## Build and lint
	golangci-lint run
	
build: fmt ## build the exe
	GOOS=windows go build -ldflags -H=windowsgui

fmt: ## format code
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

test:
	GOOS=linux richgo test ./...

build-manifest: # rebuild rsrc.syso
	# force linux tooling for rsrc (container defaults to GOOS=windows)
	GOOS=linux go install github.com/akavel/rsrc
	rsrc -manifest pick-a-browser.exe.manifest -o rsrc.syso
