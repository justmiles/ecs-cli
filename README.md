ecs-cli
===========================================
Run ad-hoc containers on ECS.

## Goals
This CLI is designed to closely resemble the `docker run` command for ECS. This will allow users to run any container against a cluster without previously having to define a task definition. It will provision the necessary AWS resources like the task definition and CloudWatch log group, start your container, and stream the logs back to stdout. Key features:

* replicate `docker run` arguments
* provision the necessary AWS resources to run a container
* the exit code from the container is mirrored by the CLI
* capture a SIGKILL event (Ctrl+c) and stop the remote container

## Usage

    Usage:
      ecs run [flags]

    Flags:
          --cluster string           ECS cluster
      -c, --count int                Spawn n tasks (default 1)
          --cpu-reservation int      CPU reservation (default 1024)
      -d, --detach                   Run the task in the background
      -e, --env stringArray          Set environment variables
          --execution-role string    Execution role ARN (required for Fargate)
          --fargate                  Launch in Fargate
      -h, --help                     help for run
      -m, --memory int               Memory limit
          --memory-reservation int   Memory reservation (default 2048)
      -n, --name string              Assign a name to the task (default "ecs-cli-app")
      -p, --publish stringArray      Publish a container's port(s) to the host
          --subnet stringArray       Subnet(s) where task should run
      -v, --volume stringArray       Map volume to ECS Container Instance

