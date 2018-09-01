VERSION=`go run main.go --version | awk '{print $$3}'`
build:
	goreleaser release --snapshot --rm-dist

tag:
	 git tag v$(VERSION) && git push --tags

release: tag
	goreleaser release --rm-dist
