package fexec

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

const (
	GET_CLUSTER_MAX_NUM = 50
)

// aws config 用インターフェース
type IFConfigService interface {
	LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error)
}

// ecs 用インターフェース
type IFEcsService interface {
	ListClusters(ctx context.Context, params *ecs.ListClustersInput, optFns ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
	ListServices(ctx context.Context, params *ecs.ListServicesInput, optFns ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
	ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
	DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
	ExecuteCommand(ctx context.Context, params *ecs.ExecuteCommandInput, optFns ...func(*ecs.Options)) (*ecs.ExecuteCommandOutput, error)
}

// ターミナル描画用インターフェース
type IFScreenDraw interface {
	ScreenDraw(options []string, label string) (string, error)
}

// aws_util 構造体
type AwsUtil struct {
	cfgif IFConfigService
	ecsif IFEcsService
	scrif IFScreenDraw
}

// 関数ラップ構造体
type WrapFnc struct{}

// LoadDefaultConfig ラッパー関数
func (t *WrapFnc) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
	return config.LoadDefaultConfig(ctx, optFns...)
}

// ScreenDraw ラッパー関数
func (t *WrapFnc) ScreenDraw(options []string, label string) (string, error) {
	return ScreenDraw(options, label)
}

// インターフェース初期化処理
func New() AwsUtil {
	awsUtil := AwsUtil{cfgif: &WrapFnc{}}
	awsUtil.scrif = &WrapFnc{}
	return awsUtil
}

// ecs クライアン初期化処理
func (awsUtil *AwsUtil) SetEcsClient(cfg aws.Config) {
	awsUtil.ecsif = ecs.NewFromConfig(cfg)
}

// バイナリファイル特定
func (awsUtil *AwsUtil) CheckBinFIle(binname string) error {
	_, err := exec.LookPath(binname)
	if err != nil {
		PrintMessage("ERR001")
	}
	return err
}

// AWS プロファイル取得
func (awsUtil *AwsUtil) FindProfile(profile string) (aws.Config, error) {
	// 利用プロファイル
	awsProfile := profile

	// 環境変数設定時にパラメータ未指定の場合
	if profile == "default" && os.Getenv("AWS_DEFAULT_PROFILE") != "" {
		// 環境変数設定時は環境変数のプロファイルを利用
		awsProfile = os.Getenv("AWS_DEFAULT_PROFILE")
	}

	// プロファイル取得
	awsCfg, err := awsUtil.cfgif.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(awsProfile),
	)
	if err != nil {
		PrintMessage("ERR002")
		return awsCfg, err
	}
	if awsCfg.Region == "" {
		PrintMessage("INF001")
	}
	return awsCfg, nil
}

// クラスター取得
func (awsUtil *AwsUtil) GetCluster() (cluster string, err error) {
	// クラスター取得パラメータ定義
	params := &ecs.ListClustersInput{
		// 取得最大件数
		MaxResults: aws.Int32(GET_CLUSTER_MAX_NUM),
	}
	// クラスター一覧取得
	resp, err := awsUtil.ecsif.ListClusters(context.TODO(), params)
	if resp == nil || err != nil {
		PrintMessage("ERR003")
		return "", err
	}
	// クラスター取得件数分ループ
	var clusters []string
	for _, v := range resp.ClusterArns {
		// クラスター名取得
		cluster := strings.Split(v, "/")[1]
		// Fargate サービス一覧取得
		services, err := awsUtil.GetServicesList(cluster)
		if err != nil {
			PrintMessage("ERR003")
			return "", err
		}
		// サービスが存在しない場合
		if services == nil {
			// 次のループへ
			continue
			// サービスが存在する場合
		} else {
			// クラスター一覧に追加
			clusters = append(clusters, cluster)
		}
	}
	// Fargate を含むクラスターが存在しない場合
	if len(clusters) <= 0 {
		PrintMessage("ERR003")
		return "", err
	}
	// クラスター一覧選択
	cluster, err = awsUtil.scrif.ScreenDraw(clusters, "cluster")
	// クラスター未選択
	if cluster == "" {
		PrintMessage("INF002")
		return "", err
	}
	return cluster, nil
}

// サービス一覧取得
func (awsUtil *AwsUtil) GetServicesList(cluster string) (services []string, err error) {
	// サービス取得パラメータ定義
	params := &ecs.ListServicesInput{
		// クラスター名
		Cluster: aws.String(cluster),
		// ローンチタイプ
		LaunchType: "FARGATE",
		// 取得最大件数
		MaxResults: aws.Int32(GET_CLUSTER_MAX_NUM),
	}
	// サービス一覧取得
	resp, err := awsUtil.ecsif.ListServices(context.TODO(), params)
	if len(resp.ServiceArns) <= 0 || err != nil {
		return nil, err
	}
	// サービス取得件数分ループ
	for _, v := range resp.ServiceArns {
		services = append(services, strings.Split(v, "/")[len(strings.Split(v, "/"))-1])
	}
	return services, nil
}

// サービス取得
func (awsUtil *AwsUtil) GetService(cluster string) (service string, err error) {
	// サービス一覧取得
	services, err := awsUtil.GetServicesList(cluster)
	if services == nil || err != nil {
		PrintMessage("ERR004")
		return "", err
	}
	// サービス一覧選択
	service, err = awsUtil.scrif.ScreenDraw(services, "service")
	// サービス未選択
	if service == "" {
		PrintMessage("INF003")
		return "", err
	}
	return service, nil
}

// タスク一覧取得
func (awsUtil *AwsUtil) GetTask(cluster string, service string) (task string, err error) {
	// タスク取得パラメータ定義
	params := &ecs.ListTasksInput{
		// クラスター名
		Cluster: aws.String(cluster),
		// ローンチタイプ
		LaunchType: "FARGATE",
		// 取得最大件数
		MaxResults: aws.Int32(GET_CLUSTER_MAX_NUM),
		// サービス名
		ServiceName: aws.String(service),
	}
	// タスク一覧取得
	resp, err := awsUtil.ecsif.ListTasks(context.TODO(), params)
	if len(resp.TaskArns) <= 0 || err != nil {
		PrintMessage("ERR005")
		return "", err
	}
	// タスク一覧を設定
	var tasks []string
	for _, v := range resp.TaskArns {
		tasks = append(tasks, strings.Split(v, "/")[len(strings.Split(v, "/"))-2]+"/"+strings.Split(v, "/")[len(strings.Split(v, "/"))-1])
	}
	// タスク一覧選択
	task, err = awsUtil.scrif.ScreenDraw(tasks, "task")
	// タスク未選択
	if task == "" {
		PrintMessage("INF004")
		return "", err
	}
	return task, nil
}

// コンテナ一覧取得
func (awsUtil *AwsUtil) GetContainer(cluster string, task string) (container string, err error) {
	taskArn := []string{task}
	// タスク定義取得パラメータ定義
	params := &ecs.DescribeTasksInput{
		// タスク ARN
		Tasks: taskArn,
		// クラスター名
		Cluster: aws.String(cluster),
	}
	// タスク詳細取得
	resp, err := awsUtil.ecsif.DescribeTasks(context.TODO(), params)
	if len(resp.Tasks[0].Containers) <= 0 || err != nil {
		PrintMessage("ERR006")
		return "", err
	}
	// コンテナ取得件数分ループ
	var containers []string
	for _, v := range resp.Tasks[0].Containers {
		containers = append(containers, *v.Name)
	}
	// コンテナ一覧選択
	container, err = awsUtil.scrif.ScreenDraw(containers, "container")
	// コンテナ未選択
	if container == "" {
		PrintMessage("INF005")
		return "", err
	}
	return container, nil
}

// execute コマンドの実行
func (awsUtil *AwsUtil) ExecuteFargate(cluster string, task string, container string) (*ecs.ExecuteCommandOutput, error) {
	// exec コマンドパラメータ
	params := &ecs.ExecuteCommandInput{
		// クラスター名
		Cluster: aws.String(cluster),
		// コマンド
		Command: aws.String("/bin/sh"),
		// コンテナ名
		Container: aws.String(container),
		// インタラクティブモード
		Interactive: true,
		// タスク ARN
		Task: aws.String(task),
	}
	// execute コマンド実行
	resp, err := awsUtil.ecsif.ExecuteCommand(context.TODO(), params)
	if resp == nil || err != nil {
		PrintMessage("INF006")
		return nil, err
	}
	return resp, nil
}
