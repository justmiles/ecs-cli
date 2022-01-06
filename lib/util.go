package ecs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
)

func buildEnvironmentKeyValuePair(environment []string) (k []*ecs.KeyValuePair) {
	if len(environment) < 1 {
		return []*ecs.KeyValuePair{}
	}
	for _, env := range environment {
		envArr := strings.SplitN(env, "=", 2)

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
	if len(publish) < 1 {
		return []*ecs.PortMapping{}
	}
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

	return k
}

func buildMountPoint(volumes []string, efsVolumes []string) (v []*ecs.Volume, k []*ecs.MountPoint) {
	if len(volumes) < 1 && len(efsVolumes) < 1 {
		return []*ecs.Volume{}, []*ecs.MountPoint{}
	}

	// Add Bind Mounts
	for i, volume := range volumes {
		av := strings.Split(volume, ":")

		// If only one port is provided, map it to both container and host
		// default to TCP map
		sourcePath := av[0]
		volumeName := "volume" + strconv.Itoa(i)

		mountPoint := ecs.MountPoint{
			ContainerPath: &sourcePath,
			SourceVolume:  aws.String(volumeName),
			ReadOnly:      aws.Bool(false),
		}

		volume := ecs.Volume{
			Name: aws.String(volumeName),
			Host: &ecs.HostVolumeProperties{
				SourcePath: aws.String(sourcePath),
			},
		}

		// Map container port, if defined
		if len(av) > 1 {
			containerPath := av[1]
			mountPoint.ContainerPath = &containerPath
		}

		// Append to the slice
		k = append(k, &mountPoint)
		v = append(v, &volume)
	}

	// Add EFS Mounts
	for i, volume := range efsVolumes {
		av := strings.Split(volume, ":")
		efsFileSystemId := av[0]
		efsDirectory := av[1]
		containerDirectory := av[2]
		volumeName := "volume-efs" + strconv.Itoa(i)

		mountPoint := ecs.MountPoint{
			ContainerPath: &containerDirectory,
			SourceVolume:  aws.String(volumeName),
			ReadOnly:      aws.Bool(false),
		}

		volume := ecs.Volume{
			Name: aws.String(volumeName),
			EfsVolumeConfiguration: &ecs.EFSVolumeConfiguration{
				// TODO parameterize this
				FileSystemId: &efsFileSystemId,
				// TODO pass in root efs mount path explicitly instead of using sourcePath
				RootDirectory: &efsDirectory,
			},
		}

		// Append to the slice
		k = append(k, &mountPoint)
		v = append(v, &volume)
	}
	return
}

// Log types
func logCloudWatchEvent(log *cloudwatchlogs.OutputLogEvent) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%v\t%v\n", yellow(time.Unix(*log.Timestamp/1000, 0)), *log.Message)
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

func getEc2InstanceIp(instanceId string) *string {
	var svc = ec2.New(sess)
	res, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{instanceId}),
	})
	logError(err)
	if res.Reservations[0].Instances[0].PublicIpAddress != nil {
		return res.Reservations[0].Instances[0].PublicIpAddress
	}
	return res.Reservations[0].Instances[0].PrivateIpAddress
}

func getSecurityGroupByName(groupName string) (*string, error) {
	var svc = ec2.New(sess)
	res, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(groupName)},
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error calling DescribeSecurityGroups %v", err)
	}

	for _, g := range res.SecurityGroups {
		return g.GroupId, nil
	}
	return nil, fmt.Errorf("unable to find security group with name %s", groupName)
}

func getSubnetsByFilter(subnetFilters []string) (subnets []*string, err error) {
	var svc = ec2.New(sess)

	var filters []*ec2.Filter

	// if filters are set, find and return subnet by filters
	for _, filter := range subnetFilters {
		s := strings.SplitN(filter, "=", 2)
		if len(s) != 2 {
			return nil, fmt.Errorf("unable to derive filter from: %s", filter)
		}
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(s[0]),
			Values: aws.StringSlice(strings.Split(s[1], ",")),
		})
	}

	result, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: filters,
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to find subnet: %s", err)
	}

	if len(result.Subnets) == 0 {
		return nil, fmt.Errorf("No subnets matching current filters")
	}

	for _, subnet := range result.Subnets {
		subnets = append(subnets, subnet.SubnetId)
	}

	return subnets, nil
}
