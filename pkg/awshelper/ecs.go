package awshelper

import (
	"context"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

const (
	maxCount = 50
)

type iFEcsService interface {
	ListClusters(ctx context.Context, params *ecs.ListClustersInput, optFns ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
	ListServices(ctx context.Context, params *ecs.ListServicesInput, optFns ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
	ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
	DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
	ExecuteCommand(ctx context.Context, params *ecs.ExecuteCommandInput, optFns ...func(*ecs.Options)) (*ecs.ExecuteCommandOutput, error)
}

type EcsService struct {
	Service iFEcsService
}

func (ecsService *EcsService) SetEcsClient(cfg aws.Config) {
	ecsService.Service = ecs.NewFromConfig(cfg)
}

func (ecsService *EcsService) GetClusters() (clusters []string, err error) {
	params := &ecs.ListClustersInput{
		MaxResults: aws.Int32(maxCount),
	}
	resp, err := ecsService.Service.ListClusters(context.TODO(), params)
	if err != nil {
		return nil, err
	}

	for _, v := range resp.ClusterArns {
		cluster := strings.Split(v, "/")[1]
		clusters = append(clusters, cluster)
	}
	sort.Strings(clusters)
	return clusters, nil
}

func (ecsService *EcsService) GetServices(cluster string) (services []string, err error) {
	params := &ecs.ListServicesInput{
		Cluster:    aws.String(cluster),
		MaxResults: aws.Int32(maxCount),
	}
	resp, err := ecsService.Service.ListServices(context.TODO(), params)
	if err != nil {
		return nil, err
	}
	for _, v := range resp.ServiceArns {
		services = append(services, strings.Split(v, "/")[len(strings.Split(v, "/"))-1])
	}
	sort.Strings(services)
	return services, nil
}

func (ecsService *EcsService) GetTasks(cluster string, service string) (tasks []string, err error) {
	params := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		MaxResults:  aws.Int32(maxCount),
		ServiceName: aws.String(service),
	}
	resp, err := ecsService.Service.ListTasks(context.TODO(), params)
	if len(resp.TaskArns) <= 0 || err != nil {
		return nil, err
	}
	for _, v := range resp.TaskArns {
		tasks = append(tasks, strings.Split(v, "/")[len(strings.Split(v, "/"))-1])
	}
	return tasks, nil
}

func (ecsService *EcsService) GetContainers(cluster string, task string) (containers []string, err error) {
	taskArn := []string{task}
	params := &ecs.DescribeTasksInput{
		Tasks:   taskArn,
		Cluster: aws.String(cluster),
	}
	resp, err := ecsService.Service.DescribeTasks(context.TODO(), params)
	if len(resp.Tasks[0].Containers) <= 0 || err != nil {
		return nil, err
	}
	for _, v := range resp.Tasks[0].Containers {
		containers = append(containers, *v.Name)
	}
	sort.Strings(containers)
	return containers, nil
}

func (ecsService *EcsService) ExecuteContainer(cluster string, task string, container string) (*ecs.ExecuteCommandOutput, error) {
	params := &ecs.ExecuteCommandInput{
		Cluster:     aws.String(cluster),
		Command:     aws.String("/bin/sh"),
		Container:   aws.String(container),
		Interactive: true,
		Task:        aws.String(task),
	}
	resp, err := ecsService.Service.ExecuteCommand(context.TODO(), params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
