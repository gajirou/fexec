package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	fexec "github.com/gajirou/fexec/util"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	profile := flag.String("p", "default", "利用プロファイル名")
	flag.Parse()

	// 初期化処理
	awsutil := fexec.New()

	// session-manager-plugin の存在チェック
	err := awsutil.CheckBinFIle("session-manager-plugin")
	if err != nil {
		return 1
	}

	var awsCfg aws.Config
	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		// aws プロファイルからコンフィグ取得
		awsCfg, err = awsutil.FindProfile(*profile)
		if err != nil {
			return 1
		}
	} else {
		// aws セッショントークンからコンフィグ取得
		awsCfg, err = awsutil.FindSessionToken()
		if err != nil {
			return 1
		}
	}

	// プロファイル未取得
	if awsCfg.Region == "" {
		return 0
	}

	// ecs クライアント設定
	awsutil.SetEcsClient(awsCfg)

	// クラスター取得
	cluster, err := awsutil.GetCluster()
	if err != nil {
		return 1
	}
	// クラスター未選択
	if cluster == "" {
		return 0
	}

	// サービス取得
	service, err := awsutil.GetService(cluster)
	if err != nil {
		return 1
	}
	// サービス未選択
	if service == "" {
		return 0
	}

	// タスク一覧取得
	task, err := awsutil.GetTask(cluster, service)
	if err != nil {
		return 1
	}
	// タスク未選択
	if task == "" {
		return 0
	}

	// コンテナ一覧取得
	container, err := awsutil.GetContainer(cluster, task)
	if container == "" || err != nil {
		return 1
	}
	// コンテナ未選択
	if container == "" {
		return 0
	}

	// execute API 実施
	execCmd, err := awsutil.ExecuteFargate(cluster, task, container)
	if err != nil {
		return 0
	}

	// ssm 用のセッション情報の余分な空白を削除
	execSes, err := json.MarshalIndent(execCmd.Session, "", "    ")
	if err != nil {
		fexec.PrintMessage("ERR999")
		panic(err)
	}

	// ssm プラグインの設定
	cmd := exec.Command("session-manager-plugin", string(execSes), awsCfg.Region, "StartSession")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	// ゴルーチンで受信したシグナルを破棄する
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	go func() {
		for {
			select {}
		}
	}()
	defer close(sig)

	// セッションを開始
	if err := cmd.Run(); err != nil {
		fexec.PrintMessage("ERR999")
		panic(err)
	}
	return 0
}
