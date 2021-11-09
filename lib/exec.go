package ecs

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ExecInput struct {
	Cluster     string
	Service     string
	Task        string
	Container   string
	Interactive bool
	Command     string
}

func GetClusters() ([]string, error) {

	clusters := []string{}

	pageNum := 0
	err := ecsClient.ListClustersPages(&ecs.ListClustersInput{},
		func(page *ecs.ListClustersOutput, lastPage bool) bool {
			pageNum++
			for _, arn := range page.ClusterArns {
				clusters = append(clusters, parseClusterName(*arn))
			}
			return pageNum <= 10
		})
	return clusters, err
}

func GetServices(cluster string) ([]string, error) {

	results := []string{}

	pageNum := 0
	err := ecsClient.ListServicesPages(&ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	},
		func(page *ecs.ListServicesOutput, lastPage bool) bool {
			pageNum++
			for _, arn := range page.ServiceArns {
				if s, err := parseServiceName(*arn); err == nil {
					results = append(results, s)
				}
			}
			return pageNum <= 10
		})
	return results, err
}

func GetRunningTasks(cluster string, service string) ([]string, error) {

	results := []string{}

	pageNum := 0
	err := ecsClient.ListTasksPages(&ecs.ListTasksInput{
		Cluster:       aws.String(cluster),
		ServiceName:   aws.String(service),
		DesiredStatus: aws.String("RUNNING"),
	},
		func(page *ecs.ListTasksOutput, lastPage bool) bool {
			pageNum++
			for _, arn := range page.TaskArns {
				// TODO get details and return more information such as uptime - start date to help idenitify tasks
				results = append(results, parseTaskId(*arn))
			}
			return pageNum <= 10
		})
	return results, err
}

func GetContainers(cluster string, task string) ([]string, error) {

	results := []string{}

	output, err := ecsClient.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   aws.StringSlice([]string{task}),
	})
	if err != nil {
		return nil, err
	}

	for _, t := range output.Tasks {
		for _, c := range t.Containers {
			results = append(results, *c.Name)
		}
	}

	return results, err
}

func ExecuteCommand(input *ExecInput) error {
	args := []string{
		"ecs",
		"execute-command",
		"--cluster",
		input.Cluster,
		"--task",
		input.Task,
		"--command",
		input.Command,
		"--container",
		input.Container,
	}
	if input.Interactive {
		args = append(args, "--interactive")
	} else {
		args = append(args, "--non-interactive")
	}

	if err := runCommand("aws", args...); err != nil {
		return err
	}
	return nil
}

func runCommand(process string, args ...string) error {
	cmd := exec.Command(process, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT)
	go func() {
		for {
			select {
			case <-sigs:
			}
		}
	}()
	defer close(sigs)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func parseClusterName(arn string) string {
	re := regexp.MustCompile("cluster/(.*?)$")
	return re.FindStringSubmatch(arn)[1]
}

func parseServiceName(arn string) (string, error) {
	re := regexp.MustCompile("service/.*/(.*?)$")
	if res := re.FindStringSubmatch(arn); len(res) > 0 {
		return res[1], nil
	} else {
		return "", fmt.Errorf("Unable to parse service name.")
	}
}

func parseTaskId(arn string) string {
	re := regexp.MustCompile("task/.*/(.*?)$")
	return re.FindStringSubmatch(arn)[1]
}
