build:
	goreleaser release --snapshot --rm-dist

tag:
	go run main.go --version | awk '{print $3}'

release: tag
	goreleaser release --rm-dist
