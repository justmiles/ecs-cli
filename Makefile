build:
	goreleaser release --snapshot --rm-dist
	
release:
	goreleaser release --rm-dist