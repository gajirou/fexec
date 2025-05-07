package awshelper_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/gajirou/fexec/pkg/awshelper"
)

type mockConfigService struct {
	cfg aws.Config
	err error
}

func (m mockConfigService) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
	return m.cfg, m.err
}

type mockEcsService struct {
	listClusterOutput    ecs.ListClustersOutput
	listServicesOutput   ecs.ListServicesOutput
	listTaskOutput       ecs.ListTasksOutput
	describeTasksOutput  ecs.DescribeTasksOutput
	executeCommandOutput ecs.ExecuteCommandOutput
	err                  error
}

func (m mockEcsService) ListClusters(ctx context.Context, params *ecs.ListClustersInput, optFns ...func(*ecs.Options)) (*ecs.ListClustersOutput, error) {
	return &m.listClusterOutput, m.err
}

func (m mockEcsService) ListServices(ctx context.Context, params *ecs.ListServicesInput, optFns ...func(*ecs.Options)) (*ecs.ListServicesOutput, error) {
	return &m.listServicesOutput, m.err
}

func (m mockEcsService) ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error) {
	return &m.listTaskOutput, m.err
}

func (m mockEcsService) DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error) {
	return &m.describeTasksOutput, m.err
}

func (m mockEcsService) ExecuteCommand(ctx context.Context, params *ecs.ExecuteCommandInput, optFns ...func(*ecs.Options)) (*ecs.ExecuteCommandOutput, error) {
	return &m.executeCommandOutput, m.err
}

func TestFindProfile(t *testing.T) {
	cases := []struct {
		name              string
		profile           string
		defaultProfileEnv bool
		sessionTokenEnv   bool
		mockError         error
		errFlag           bool
	}{
		{
			name:              "正常パターン:AWS_SESSION_TOKEN利用 ",
			profile:           "default",
			defaultProfileEnv: false,
			sessionTokenEnv:   true,
			mockError:         nil,
		},
		{
			name:              "正常パターン:AWS_SESSION_TOKEN利用+プロファイル設定有り",
			profile:           "default",
			defaultProfileEnv: true,
			sessionTokenEnv:   true,
			mockError:         nil,
		},
		{
			name:              "正常パターン:defaultプロファイル設定",
			profile:           "default",
			defaultProfileEnv: true,
			sessionTokenEnv:   false,
			mockError:         nil,
		},
		{
			name:              "正常パターン:プロファイル指定",
			profile:           "profile",
			defaultProfileEnv: true,
			sessionTokenEnv:   false,
			mockError:         nil,
		},
		{
			name:              "異常パターン:AWS_SESSION_TOKEN利用 ",
			profile:           "default",
			defaultProfileEnv: false,
			sessionTokenEnv:   true,
			mockError:         errors.New("error"),
		},
		{
			name:              "異常パターン:AWS_SESSION_TOKEN利用+プロファイル設定有り",
			profile:           "default",
			defaultProfileEnv: true,
			sessionTokenEnv:   true,
			mockError:         errors.New("error"),
		},
		{
			name:              "異常パターン:defaultプロファイル設定",
			profile:           "default",
			defaultProfileEnv: true,
			sessionTokenEnv:   false,
			mockError:         errors.New("error"),
		},
		{
			name:              "異常パターン:プロファイル指定",
			profile:           "profile",
			defaultProfileEnv: true,
			sessionTokenEnv:   false,
			mockError:         errors.New("error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockConfigService := &mockConfigService{err: c.mockError}
			mockService := awshelper.ConfigService{Service: mockConfigService}

			if c.defaultProfileEnv {
				t.Setenv("AWS_DEFAULT_PROFILE", "profile")
			}
			if c.sessionTokenEnv {
				t.Setenv("AWS_SESSION_TOKEN", "sessiontoken")
			}

			_, err := mockService.FindAWSCredential(c.profile)
			if c.mockError != nil {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値に予期せぬエラーが含まれています。")
				}
			}
		})
	}
}

func TestGetClusters(t *testing.T) {
	cases := []struct {
		name      string
		resp      ecs.ListClustersOutput
		mockError error
	}{
		{
			name: "正常パターン",
			resp: ecs.ListClustersOutput{
				ClusterArns: []string{
					"arn:aws:ecs:ap-northeast-1:111111111111:cluster/cluster1",
					"arn:aws:ecs:ap-northeast-1:111111111111:cluster/cluster2",
				},
			},
			mockError: nil,
		},
		{
			name:      "異常パターン",
			resp:      ecs.ListClustersOutput{},
			mockError: errors.New("error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockEcsService := &mockEcsService{listClusterOutput: c.resp, err: c.mockError}
			mockService := awshelper.EcsService{Service: mockEcsService}

			_, err := mockService.GetClusters()
			if c.mockError != nil {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値に予期せぬエラーが含まれています。")
				}
			}
		})
	}
}

func TestGetServices(t *testing.T) {
	cases := []struct {
		name      string
		resp      ecs.ListServicesOutput
		mockError error
	}{
		{
			name: "正常パターン",
			resp: ecs.ListServicesOutput{
				ServiceArns: []string{
					"arn:aws:ecs:ap-northeast-1:111111111111:cluster/cluster1/service1",
					"arn:aws:ecs:ap-northeast-1:111111111111:cluster/cluster1/service2",
				},
			},
			mockError: nil,
		},
		{
			name:      "異常パターン",
			resp:      ecs.ListServicesOutput{},
			mockError: errors.New("error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockEcsService := &mockEcsService{listServicesOutput: c.resp, err: c.mockError}
			mockService := awshelper.EcsService{Service: mockEcsService}

			_, err := mockService.GetServices("cluster")
			if c.mockError != nil {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値に予期せぬエラーが含まれています。")
				}
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	cases := []struct {
		name      string
		resp      ecs.ListTasksOutput
		mockError error
	}{
		{
			name: "正常パターン",
			resp: ecs.ListTasksOutput{
				TaskArns: []string{
					"arn:aws:ecs:ap-northeast-1:cluster/task1/task1-1",
					"arn:aws:ecs:ap-northeast-1:cluster/task1/task1-2",
				},
			},
			mockError: nil,
		},
		{
			name:      "異常パターン:通常エラー",
			resp:      ecs.ListTasksOutput{},
			mockError: errors.New("error"),
		},
		{
			name:      "異常パターン:サービスに紐づくタスクが存在しない",
			resp:      ecs.ListTasksOutput{},
			mockError: nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockEcsService := &mockEcsService{listTaskOutput: c.resp, err: c.mockError}
			mockService := awshelper.EcsService{Service: mockEcsService}

			_, err := mockService.GetTasks("cluster", "service")
			if c.mockError != nil {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値に予期せぬエラーが含まれています。")
				}
			}
		})
	}
}

func TestGetContainers(t *testing.T) {
	cases := []struct {
		name      string
		resp      ecs.DescribeTasksOutput
		mockError error
	}{
		{
			name: "正常パターン",
			resp: ecs.DescribeTasksOutput{
				Tasks: []types.Task{
					{
						Containers: []types.Container{
							{
								Name: aws.String("conteinar"),
							},
						},
					},
				},
			},
			mockError: nil,
		},
		{
			name: "異常パターン:通常エラー",
			resp: ecs.DescribeTasksOutput{
				Tasks: []types.Task{
					{
						Containers: []types.Container{
							{
								Name: aws.String("conteinar"),
							},
						},
					},
				},
			},
			mockError: errors.New("error"),
		},
		{
			name: "異常パターン:タスクに紐づくコンテナが存在しない",
			resp: ecs.DescribeTasksOutput{
				Tasks: []types.Task{{}},
			},
			mockError: nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockEcsService := &mockEcsService{describeTasksOutput: c.resp, err: c.mockError}
			mockService := awshelper.EcsService{Service: mockEcsService}

			_, err := mockService.GetContainers("cluster", "service")
			if c.mockError != nil {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値に予期せぬエラーが含まれています。")
				}
			}
		})
	}
}

func TestExecuteContainer(t *testing.T) {
	cases := []struct {
		name      string
		resp      ecs.ExecuteCommandOutput
		mockError error
	}{
		{
			name: "正常パターン",
			resp: ecs.ExecuteCommandOutput{
				Session: &types.Session{
					SessionId:  aws.String("session"),
					StreamUrl:  aws.String("url"),
					TokenValue: aws.String("token"),
				},
			},
			mockError: nil,
		},
		{
			name:      "異常パターン",
			resp:      ecs.ExecuteCommandOutput{},
			mockError: errors.New("error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockEcsService := &mockEcsService{executeCommandOutput: c.resp, err: c.mockError}
			mockService := awshelper.EcsService{Service: mockEcsService}

			_, err := mockService.ExecuteContainer("cluster", "task", "container")
			if c.mockError != nil {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値に予期せぬエラーが含まれています。")
				}
			}
		})
	}
}
