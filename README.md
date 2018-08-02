> $ ecs run hello-world

    Hello from Docker!
    This message shows that your installation appears to be working correctly.

    To generate this message, Docker took the following steps:
     1. The Docker client contacted the Docker daemon.
     2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
        (amd64)
     3. The Docker daemon created a new container from that image which runs the
        executable that produces the output you are currently reading.
     4. The Docker daemon streamed that output to the Docker client, which sent it
        to your terminal.

    To try something more ambitious, you can run an Ubuntu container with:
    
     $ docker run -it ubuntu bash

    Share images, automate workflows, and more with a free Docker ID:
     https://hub.docker.com/

    For more examples and ideas, visit:
     https://docs.docker.com/engine/userguide/

$ ecs run --detach hello-world

    service xxx started.

$ ecs run --restart hello-world
    
    service xxx started.


## Usage

    ECS Options:
      --cluster                                                           [required]
      --region                                                  (default: us-east-1)
      --fargate                                                 (default: us-east-1)

    RUN Options:
      -n, --name                           Assign a name to the task
      -d, --detach                         Run the task in the background
      -c, --count                          Spawn x tasks
      -e, --env   list                     Set environment variables
      -h, --hostname   string               Set environment variables
      -p, --publish   list                  Publish a container's port(s) to the host
          --public   string                 Assign public IP
          --security-groups   list          Attach security groups to task
      -v, --volume list                     
      -m, --memory bytes                     Memory limit
          --memory-reservation bytes         Memory reservation