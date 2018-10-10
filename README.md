## Usage

    Usage:
      ecs run [flags]

    Flags:
          --cluster string                ECS cluster
      -c, --count int                     Spawn n tasks (default 1)
          --cpu-reservation int           CPU reservation (default 1024)
      -d, --detach                        [TODO] Run the task in the background
      -e, --env stringArray               Set environment variables
          --execution-role string         Execution role ARN (required for Fargate)
          --fargate                       Launch in Fargate
      -h, --help                          help for run
      -m, --memory int                    Memory limit
          --memory-reservation int        Memory reservation (default 2048)
      -n, --name string                   Assign a name to the task (default "ecs-cli-app")
          --public                        [TODO] Assign public IP
      -p, --publish stringArray           Publish a container's port(s) to the host
          --security-groups stringArray   [TODO] Attach security groups to task
          --subnet stringArray            Subnet(s) where task should run
      -v, --volume stringArray            Map volume to ECS Container Instance
