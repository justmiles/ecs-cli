build:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release --rm-dist

task-def:
	go run main.go run-task-def --cluster qa --security-groups <security-group-name> --subnet-filter "tag:Name=<subnet-name>" --family <task-def-name> --image-version latest

docker-build:
	docker buildx build --platform=linux/amd64 -t justmiles/ecs-cli . 

docker-run:
	docker run --platform=linux/amd64 -it -e AWS_SESSION_TOKEN -e AWS_SECRET_ACCESS_KEY -e AWS_ACCESS_KEY_ID -e AWS_DEFAULT_REGION=us-east-1 justmiles/ecs-cli exec

docker-build-slim:
	docker buildx build --platform=linux/amd64 -f slim.Dockerfile -t justmiles/ecs-cli:slim . 

docker-run-slim:
	docker run --platform=linux/amd64 -it -e AWS_SESSION_TOKEN -e AWS_SECRET_ACCESS_KEY -e AWS_ACCESS_KEY_ID -e AWS_DEFAULT_REGION=us-east-1 justmiles/ecs-cli:slim exec