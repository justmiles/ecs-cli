ecs-cli
===========================================
Run ad-hoc containers on ECS.

This does not aim to replace [the offical AWS ECS CLI](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_CLI.html), but complement it. The official CLI doesn't support `docker run` like commands. 

## Goals
This CLI is designed to closely resemble the `docker run` command for ECS. This will allow users to run any container against a cluster without previously having to define a task definition. It will provision the necessary AWS resources like the task definition and CloudWatch log group, start your container, and stream the logs back to stdout.

* [x] replicate `docker run` arguments
* [x] provision the necessary AWS resources to run a container
    * [ ] allow custom task-definition name
* [x] the exit code from the container is mirrored by the CLI
* [x] capture a SIGKILL event (Ctrl+c) and stop the remote container
* [x] replicate `docker exec` arguments using SSM agent
* [x] support executing an existing task definition with optional image version bump
* [ ] Only create new task definition if it differs from the last active task defition
* [ ] Specify (or maybe create?) cluster capacity provider for fargate spot support

## Installation
[Download the latest](https://github.com/justmiles/ecs-cli/releases) release for your OS and extract to your PATH.

**Homebrew**
```bash
brew install chrispruitt/tap/ecs
```

## Example

```
➜  ~ ecs run --cluster ops bash ping -c 5 google.com
Creating task definition
Running task definition: arn:aws:ecs:us-east-1:000000000000:task-definition/ephemeral-task-from-ecs-cli:1
Streaming from Cloudwatch Logs
https://console.aws.amazon.com/ecs/home?#/clusters/ops/tasks/00000000000000000000000000000000/details
Container is starting on EC2 instance i-00000000000000000 (10.100.1.25).
2019-01-19 20:25:08 -0600 CST	PING google.com (172.217.7.238): 56 data bytes
2019-01-19 20:25:08 -0600 CST	64 bytes from 172.217.7.238: seq=0 ttl=46 time=2.020 ms
2019-01-19 20:25:09 -0600 CST	64 bytes from 172.217.7.238: seq=1 ttl=46 time=1.303 ms
2019-01-19 20:25:10 -0600 CST	64 bytes from 172.217.7.238: seq=2 ttl=46 time=1.263 ms
2019-01-19 20:25:11 -0600 CST	64 bytes from 172.217.7.238: seq=3 ttl=46 time=1.230 ms
2019-01-19 20:25:12 -0600 CST	64 bytes from 172.217.7.238: seq=4 ttl=46 time=1.313 ms
2019-01-19 20:25:12 -0600 CST	--- google.com ping statistics ---
2019-01-19 20:25:12 -0600 CST	5 packets transmitted, 5 packets received, 0% packet loss
2019-01-19 20:25:12 -0600 CST	round-trip min/avg/max = 1.230/1.425/2.020 ms
Task arn:aws:ecs:us-east-1:000000000000:task/ops/00000000000000000000000000000000 has stopped (exit code 0):
	Essential container in task exited
All containers have exited
```

## Usage

    Usage:
    ecs run [flags]

    Flags:
        --cli-role string               An IAM role ARN to assume before creating/executing a task
        --cluster string                ECS cluster
    -c, --count int                     Spawn n tasks (default 1)
        --cpu-reservation int           CPU reservation (default 256)
        --debug                         Verbose logging
    -d, --detach                        Run the task in the background
        --efs-volume stringArray        Map EFS volume to ECS Container Instance (ex. fs-23kj2f:/efs/dir:/container/mnt/dir)
    -e, --env stringArray               Set environment variables
        --execution-role string         Execution role ARN (required for Fargate)
        --family string                 Family for ECS task
        --fargate                       Launch in Fargate
    -h, --help                          help for run
    -m, --memory int                    Memory limit
        --memory-reservation int        Memory reservation (default 2048)
    -n, --name string                   Assign a name to the task (default "ephemeral-task-from-ecs-cli")
        --no-deregister                 do not deregister the task definition
        --public                        assign public ip
    -p, --publish stringArray           Publish a container's port(s) to the host
        --role string                   Task role ARN
        --security-groups stringArray   attach security groups to task
        --subnet-filter stringArray     'Key=Value' filters for your subnet, eg tag:Name=private
    -t, --tag stringArray               Tag task definition on creation (eg key=value). Multiple uses for multiple tags
    -v, --volume stringArray            Map volume to ECS Container Instance


## Note

The slim docker image is much smaller, but does not support the exec command.
