build:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release --rm-dist

task-def:
	go run main.go run-task-def --cluster qa --security-groups <security-group-name> --subnet-filter "tag:Name=<subnet-name>" --family <task-def-name> --image-version latest