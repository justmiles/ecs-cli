package ecs

import "github.com/aws/aws-sdk-go/service/ecs"
import "github.com/aws/aws-sdk-go/aws"
import "fmt"
import "os"
import "strings"
import "strconv"
import "time"
import "github.com/fatih/color"
import "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

func buildEnvironmentKeyValuePair(environment []string) (k []*ecs.KeyValuePair) {

	for _, env := range environment {
		envArr := strings.Split(env, "=")

		// Recieved a --env MYVAR=test
		if len(envArr) > 1 {
			k = append(k, &ecs.KeyValuePair{
				Name:  &envArr[0],
				Value: &envArr[1],
			})

			// Recieved a --env MYVAR
			// Try to get a value from MYVAR environment variable
		} else {
			if value, ok := os.LookupEnv(envArr[0]); ok {
				k = append(k, &ecs.KeyValuePair{
					Name:  &envArr[0],
					Value: &value,
				})
			}
		}
	}
	return
}

func buildPortMapping(publish []string) (k []*ecs.PortMapping) {

	for _, env := range publish {
		envArr := strings.Split(env, ":")

		// If only one port is provided, map it to both container and host
		// default to TCP map
		hostPort, _ := strconv.ParseInt(envArr[0], 10, 64)
		portMap := ecs.PortMapping{
			ContainerPort: &hostPort,
			HostPort:      &hostPort,
			Protocol:      aws.String("tcp"),
		}

		// Map container port, if defined
		if len(envArr) > 1 {
			containerPort, _ := strconv.ParseInt(envArr[1], 10, 64)
			portMap.ContainerPort = &containerPort
		}

		// Override tcp protocol with whatever is provided (Should just be "udp")
		if len(envArr) > 2 {
			portMap.Protocol = aws.String(envArr[3])
		}

		// Append to the slice
		k = append(k, &portMap)
	}

	return
}

// Log types
func logCloudWatchEvent(log *cloudwatchlogs.OutputLogEvent) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("[%v]\t%v\n", yellow(time.Unix(*log.Timestamp/1000, 0)), *log.Message)
}

func logInfo(s string) {
	color.Green(s)
}

func logFatal(s string) {
	color.Red(s)
	os.Exit(1)
}

func logFatalError(e error) {
	if e != nil {
		color.Red(e.Error())
		os.Exit(1)
	}
}

func logError(e error) {
	if e != nil {
		color.Red(e.Error())
	}
}

func logWarning(s string) {
	color.Yellow(s)
}
