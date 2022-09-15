package fexec

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// 標準出力キャプチャ構造体
type PrintCapcha struct {
	stdout *os.File
	r      *os.File
	w      *os.File
	output string
}

// 標準出力キャプチャ前処理
func (p *PrintCapcha) PreCapPrint(t *testing.T) {
	t.Helper()

	p.stdout = os.Stdout
	p.r, p.w, _ = os.Pipe()
	os.Stdout = p.w
}

// 標準出力キャプチャ後処理
func (p *PrintCapcha) PostCapPrint(t *testing.T) {
	t.Helper()

	p.w.Close()
	os.Stdout = p.stdout
	buff := bytes.Buffer{}
	io.Copy(&buff, p.r)
	p.output = strings.TrimRight(buff.String(), "\n")
}

// mock 構造体
type AwsUtilMock struct{}

func (a AwsUtilMock) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
	switch os.Getenv("CASE") {
	case "case001":
		cfg.Region = "ap-northeast-1"
		err = nil
	case "case002":
		cfg.Region = "ap-northeast-1"
		err = nil
	case "case003":
		cfg.Region = "ap-northeast-1"
		err = nil
	case "case004":
		cfg.Region = "ap-northeast-1"
		err = nil
	case "case005":
		cfg.Region = "ap-northeast-1"
		err = errors.New("err")
	case "case006":
		cfg.Region = "ap-northeast-1"
		err = errors.New("err")
	case "case007":
		cfg.Region = ""
		err = errors.New("err")
	case "case008":
		cfg.Region = ""
		err = errors.New("err")
	case "case009":
		cfg.Region = ""
		err = nil
	case "case010":
		cfg.Region = ""
		err = nil
	case "case011":
		cfg.Region = ""
		err = nil
	case "case012":
		cfg.Region = ""
		err = nil
	default:
		cfg.Region = ""
		err = nil
	}
	return cfg, err
}

func (a AwsUtilMock) ListClusters(ctx context.Context, params *ecs.ListClustersInput, optFns ...func(*ecs.Options)) (resp *ecs.ListClustersOutput, err error) {
	switch os.Getenv("CASE") {
	case "case101":
		resp = nil
		err = errors.New("err")
	case "case102":
		resp = nil
		err = nil
	case "case103":
		resp = &ecs.ListClustersOutput{
			ClusterArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case103",
			},
		}
		err = nil
	case "case104":
		resp = &ecs.ListClustersOutput{
			ClusterArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case104",
			},
		}
		err = nil
	case "case105":
		resp = &ecs.ListClustersOutput{
			ClusterArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case105",
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case105-test",
			},
		}
		err = nil
	case "case106":
		resp = &ecs.ListClustersOutput{
			ClusterArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case106",
			},
		}
		err = nil
	default:
		resp = nil
		err = nil
	}
	return resp, err
}

func (a AwsUtilMock) ListServices(ctx context.Context, params *ecs.ListServicesInput, optFns ...func(*ecs.Options)) (resp *ecs.ListServicesOutput, err error) {
	switch os.Getenv("CASE") {
	case "case103":
		resp = &ecs.ListServicesOutput{
			ServiceArns: []string{},
		}
		err = nil
	case "case104":
		resp = &ecs.ListServicesOutput{
			ServiceArns: []string{},
		}
		err = errors.New("err")
	case "case105":
		resp = &ecs.ListServicesOutput{
			ServiceArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case105/case105",
			},
		}
		err = nil
	case "case106":
		resp = &ecs.ListServicesOutput{
			ServiceArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case106/case106",
			},
		}
		err = nil
	case "case201":
		resp = &ecs.ListServicesOutput{
			ServiceArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case201/case201",
			},
		}
		err = errors.New("err")
	case "case202":
		resp = &ecs.ListServicesOutput{
			ServiceArns: []string{},
		}
		err = nil
	case "case203":
		resp = &ecs.ListServicesOutput{
			ServiceArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case203/case203",
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case203/case203-test",
			},
		}
		err = nil
	case "case204":
		resp = &ecs.ListServicesOutput{
			ServiceArns: []string{
				"arn:aws:ecs:ap-northeast-1:111111111111:cluster/case204/case204",
			},
		}
		err = nil
	default:
		resp = nil
		err = nil
	}
	return resp, err
}

func (a AwsUtilMock) ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (resp *ecs.ListTasksOutput, err error) {
	switch os.Getenv("CASE") {
	case "case301":
		resp = &ecs.ListTasksOutput{
			TaskArns: []string{
				"arn:aws:ecs:ap-northeast-1:task/case301/11111111111111111111111",
			},
		}
		err = errors.New("err")
	case "case302":
		resp = &ecs.ListTasksOutput{
			TaskArns: []string{},
		}
		err = nil
	case "case303":
		resp = &ecs.ListTasksOutput{
			TaskArns: []string{
				"arn:aws:ecs:ap-northeast-1:task/case303/11111111111111111111111",
				"arn:aws:ecs:ap-northeast-1:task/case303/22222222222222222222222",
			},
		}
		err = nil
	case "case304":
		resp = &ecs.ListTasksOutput{
			TaskArns: []string{
				"arn:aws:ecs:ap-northeast-1:task/case304/11111111111111111111111",
			},
		}
		err = nil
	}
	return resp, err
}

func (a AwsUtilMock) DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (resp *ecs.DescribeTasksOutput, err error) {
	switch os.Getenv("CASE") {
	case "case401":
		resp = &ecs.DescribeTasksOutput{
			Tasks: []types.Task{
				{
					Containers: []types.Container{
						{
							Name: aws.String("case401"),
						},
					},
				},
			},
		}
		err = errors.New("err")
	case "case402":
		resp = &ecs.DescribeTasksOutput{
			Tasks: []types.Task{
				{},
			},
		}
		err = nil
	case "case403":
		resp = &ecs.DescribeTasksOutput{
			Tasks: []types.Task{
				{
					Containers: []types.Container{
						{
							Name: aws.String("case403"),
						},
						{
							Name: aws.String("case403-container"),
						},
					},
				},
			},
		}
		err = nil
	case "case404":
		resp = &ecs.DescribeTasksOutput{
			Tasks: []types.Task{
				{
					Containers: []types.Container{
						{
							Name: aws.String("case404"),
						},
					},
				},
			},
		}
		err = nil
	default:
		resp = nil
		err = nil
	}
	return resp, err
}

func (a AwsUtilMock) ExecuteCommand(ctx context.Context, params *ecs.ExecuteCommandInput, optFns ...func(*ecs.Options)) (resp *ecs.ExecuteCommandOutput, err error) {
	switch os.Getenv("CASE") {
	case "case501":
		resp = nil
		err = errors.New("err")
	case "case502":
		resp = &ecs.ExecuteCommandOutput{
			Session: &types.Session{
				SessionId:  aws.String("case502"),
				StreamUrl:  aws.String("case502"),
				TokenValue: aws.String("case502"),
			},
		}
		err = nil
	default:
		resp = nil
		err = nil
	}
	return resp, err
}

func (a AwsUtilMock) ScreenDraw(options []string, label string) (ask string, err error) {
	switch os.Getenv("CASE") {
	case "case105", "case203", "case303", "case403":
		ask = options[0]
		err = nil
	default:
		ask = ""
		err = nil
	}
	return ask, err
}

func TestCheckBinFIle(t *testing.T) {
	// テストケース
	cases := []struct {
		// ケース名
		name string
		// 出力メッセージ
		message string
		// 引数
		arg string
		// エラーケースフラグ
		err bool
	}{
		{
			name:    "case ファイル特定",
			message: "",
			arg:     "ls",
			err:     false,
		},
		{
			name:    "case ファイル未特定",
			message: GetMessage("ERR001"),
			arg:     "lss",
			err:     true,
		},
	}
	// aws_util 構造体定義
	a := AwsUtil{}
	// キャプチャ用構造体定義
	p := PrintCapcha{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// メッセージラベルが設定されていればキャプチャ前処理実行
			if c.message != "" {
				p.PreCapPrint(t)
			}
			err := a.CheckBinFIle(c.arg)
			// メッセージラベルが設定されていればキャプチャ後処理実行
			if c.message != "" {
				p.PostCapPrint(t)
				expect := c.message
				if p.output != expect {
					t.Errorf("期待するメッセージと異なります。\n出力メッセージ - \n%v\n期待するメッセージ - \n%v", p.output, expect)
				}
			}
			// c.err フラグでエラー処理のテストを分岐
			if c.err {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値にエラーが含まれています。")
				}
			}
		})
	}
}

func TestFindProfile(t *testing.T) {
	cases := []struct {
		name    string
		message string
		fcase   string
		param   string
		env     string
		err     bool
	}{
		{
			name:    "case パラメータ・環境変数指定プロファイル取得",
			message: "",
			fcase:   "case001",
			param:   "default",
			env:     "default",
			err:     false,
		},
		{
			name:    "case パラメータ指定・環境変数未指定プロファイル取得",
			message: "",
			fcase:   "case002",
			param:   "default",
			env:     "",
			err:     false,
		},
		{
			name:    "case パラメータ未指定・環境変数指定プロファイル取得",
			message: "",
			fcase:   "case003",
			param:   "",
			env:     "default",
			err:     false,
		},
		{
			name:    "case パラメータ・環境変数未指定プロファイル取得",
			message: "",
			fcase:   "case004",
			param:   "",
			env:     "",
			err:     false,
		},
		{
			name:    "case パラメータ・環境変数指定エラー発生",
			message: GetMessage("ERR002"),
			fcase:   "case005",
			param:   "default",
			env:     "default",
			err:     true,
		},
		{
			name:    "case パラメータ指定・環境変数未指定エラー発生",
			message: GetMessage("ERR002"),
			fcase:   "case006",
			param:   "default",
			env:     "",
			err:     true,
		},
		{
			name:    "case パラメータ未指定・環境変数指定エラー発生",
			message: GetMessage("ERR002"),
			fcase:   "case007",
			param:   "",
			env:     "default",
			err:     true,
		},
		{
			name:    "case パラメータ・環境変数未指定エラー発生",
			message: GetMessage("ERR002"),
			fcase:   "case008",
			param:   "",
			env:     "",
			err:     true,
		},
		{
			name:    "case パラメータ・環境変数指定プロファイルなし",
			message: GetMessage("INF001"),
			fcase:   "case009",
			param:   "default",
			env:     "default",
			err:     false,
		},
		{
			name:    "case パラメータ指定・環境変数未指定プロファイルなし",
			message: GetMessage("INF001"),
			fcase:   "case010",
			param:   "default",
			env:     "",
			err:     false,
		},
		{
			name:    "case パラメータ未指定・環境変数指定プロファイルなし",
			message: GetMessage("INF001"),
			fcase:   "case011",
			param:   "",
			env:     "default",
			err:     false,
		},
		{
			name:    "case パラメータ・環境変数未指定プロファイルなし",
			message: GetMessage("INF001"),
			fcase:   "case012",
			param:   "",
			env:     "",
			err:     false,
		},
	}
	a := AwsUtil{cfgif: AwsUtilMock{}}
	p := PrintCapcha{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("CASE", c.fcase)
			if c.env != "" {
				t.Setenv("AWS_DEFAULT_PROFILE", c.env)
			}
			if c.message != "" {
				p.PreCapPrint(t)
			}
			_, err := a.FindProfile(c.param)
			if c.message != "" {
				p.PostCapPrint(t)
				expect := c.message
				if p.output != expect {
					t.Errorf("期待するメッセージと異なります。\n出力メッセージ - \n%v\n期待するメッセージ - \n%v", p.output, expect)
				}
			}
			if c.err {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値にエラーが含まれています。")
				}
			}
		})
	}
}

func TestGetCluster(t *testing.T) {
	cases := []struct {
		name    string
		message string
		fcase   string
		err     bool
	}{
		{
			name:    "case クラスター一覧取得時エラー発生",
			message: GetMessage("ERR003"),
			fcase:   "case101",
			err:     true,
		},
		{
			name:    "case クラスター未検出",
			message: GetMessage("ERR003"),
			fcase:   "case102",
			err:     false,
		},
		{
			name:    "case クラスターに紐づくFargateサービスなし",
			message: GetMessage("ERR003"),
			fcase:   "case103",
			err:     false,
		},
		{
			name:    "case クラスターに紐づくサービス取得時エラー発生",
			message: GetMessage("ERR003"),
			fcase:   "case104",
			err:     true,
		},
		{
			name:    "case クラスター取得",
			message: "",
			fcase:   "case105",
			err:     false,
		},
		{
			name:    "case クラスター未選択",
			message: GetMessage("INF002"),
			fcase:   "case106",
			err:     false,
		},
	}
	a := AwsUtil{
		ecsif: AwsUtilMock{},
		scrif: AwsUtilMock{},
	}
	p := PrintCapcha{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("CASE", c.fcase)
			if c.message != "" {
				p.PreCapPrint(t)
			}
			rtn, err := a.GetCluster()
			if rtn != "" {
				if rtn != c.fcase {
					t.Errorf("期待する返り値と異なります。\n返り値 - \n%v\n期待する返り値 - \n%v", rtn, c.fcase)
				}
			}
			if c.message != "" {
				p.PostCapPrint(t)
				expect := c.message
				if p.output != expect {
					t.Errorf("期待するメッセージと異なります。\n出力メッセージ - \n%v\n期待するメッセージ - \n%v", p.output, expect)
				}
			}
			if c.err {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値にエラーが含まれています。")
				}
			}
		})
	}
}

func TestGetService(t *testing.T) {
	cases := []struct {
		name    string
		message string
		fcase   string
		err     bool
	}{
		{
			name:    "case サービス一覧取得時エラー発生",
			message: GetMessage("ERR004"),
			fcase:   "case201",
			err:     true,
		},
		{
			name:    "case サービス未検出",
			message: GetMessage("ERR004"),
			fcase:   "case202",
			err:     false,
		},
		{
			name:    "case サービス取得",
			message: "",
			fcase:   "case203",
			err:     false,
		},
		{
			name:    "case サービス未選択",
			message: GetMessage("INF003"),
			fcase:   "case204",
			err:     false,
		},
	}
	a := AwsUtil{
		ecsif: AwsUtilMock{},
		scrif: AwsUtilMock{},
	}
	p := PrintCapcha{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("CASE", c.fcase)
			if c.message != "" {
				p.PreCapPrint(t)
			}
			rtn, err := a.GetService("")
			if rtn != "" {
				if rtn != c.fcase {
					t.Errorf("期待する返り値と異なります。\n返り値 - \n%v\n期待する返り値 - \n%v", rtn, c.fcase)
				}
			}
			if c.message != "" {
				p.PostCapPrint(t)
				expect := c.message
				if p.output != expect {
					t.Errorf("期待するメッセージと異なります。\n出力メッセージ - \n%v\n期待するメッセージ - \n%v", p.output, expect)
				}
			}
			if c.err {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値にエラーが含まれています。")
				}
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	cases := []struct {
		name    string
		message string
		fcase   string
		err     bool
	}{
		{
			name:    "case タスク一覧取得時エラー発生",
			message: GetMessage("ERR005"),
			fcase:   "case301",
			err:     true,
		},
		{
			name:    "case タスク未検出",
			message: GetMessage("ERR005"),
			fcase:   "case302",
			err:     false,
		},
		{
			name:    "case タスク取得",
			message: "",
			fcase:   "case303",
			err:     false,
		},
		{
			name:    "case タスク未選択",
			message: GetMessage("INF004"),
			fcase:   "case304",
			err:     false,
		},
	}
	a := AwsUtil{
		ecsif: AwsUtilMock{},
		scrif: AwsUtilMock{},
	}
	p := PrintCapcha{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("CASE", c.fcase)
			if c.message != "" {
				p.PreCapPrint(t)
			}
			rtn, err := a.GetTask("", "")
			if rtn != "" {
				expect := "case303/11111111111111111111111"
				if rtn != expect {
					t.Errorf("期待する返り値と異なります。\n返り値 - \n%v\n期待する返り値 - \n%v", rtn, expect)
				}
			}
			if c.message != "" {
				p.PostCapPrint(t)
				expect := c.message
				if p.output != expect {
					t.Errorf("期待するメッセージと異なります。\n出力メッセージ - \n%v\n期待するメッセージ - \n%v", p.output, expect)
				}
			}
			if c.err {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値にエラーが含まれています。")
				}
			}
		})
	}
}

func TestGetContainer(t *testing.T) {
	cases := []struct {
		name    string
		message string
		fcase   string
		err     bool
	}{
		{
			name:    "case コンテナ一覧取得時エラー発生",
			message: GetMessage("ERR006"),
			fcase:   "case401",
			err:     true,
		},
		{
			name:    "case コンテナ未検出",
			message: GetMessage("ERR006"),
			fcase:   "case402",
			err:     false,
		},
		{
			name:    "case コンテナ取得",
			message: "",
			fcase:   "case403",
			err:     false,
		},
		{
			name:    "case コンテナ未選択",
			message: GetMessage("INF005"),
			fcase:   "case404",
			err:     false,
		},
	}
	a := AwsUtil{
		ecsif: AwsUtilMock{},
		scrif: AwsUtilMock{},
	}
	p := PrintCapcha{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("CASE", c.fcase)
			if c.message != "" {
				p.PreCapPrint(t)
			}
			rtn, err := a.GetContainer("", "")
			if rtn != "" {
				if rtn != c.fcase {
					t.Errorf("期待する返り値と異なります。\n返り値 - \n%v\n期待する返り値 - \n%v", rtn, c.fcase)
				}
			}
			if c.message != "" {
				p.PostCapPrint(t)
				expect := c.message
				if p.output != expect {
					t.Errorf("期待するメッセージと異なります。\n出力メッセージ - \n%v\n期待するメッセージ - \n%v", p.output, expect)
				}
			}
			if c.err {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値にエラーが含まれています。")
				}
			}
		})
	}
}

func TestExecuteFargate(t *testing.T) {
	cases := []struct {
		name    string
		message string
		fcase   string
		err     bool
	}{
		{
			name:    "case execute command エラー発生",
			message: GetMessage("INF006"),
			fcase:   "case501",
			err:     true,
		},
		{
			name:    "case execute command レスポンス取得",
			message: "",
			fcase:   "case502",
			err:     false,
		},
	}
	a := AwsUtil{
		ecsif: AwsUtilMock{},
	}
	p := PrintCapcha{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("CASE", c.fcase)
			if c.message != "" {
				p.PreCapPrint(t)
			}
			rtn, err := a.ExecuteFargate("", "", "")
			if rtn != nil {
				if *rtn.Session.SessionId != c.fcase {
					t.Errorf("期待する返り値と異なります。\n返り値 - \n%v\n期待する返り値 - \n%v", *rtn.Session.SessionId, c.fcase)
				}
			}
			if c.message != "" {
				p.PostCapPrint(t)
				expect := c.message
				if p.output != expect {
					t.Errorf("期待するメッセージと異なります。\n出力メッセージ - \n%v\n期待するメッセージ - \n%v", p.output, expect)
				}
			}
			if c.err {
				if err == nil {
					t.Error("関数の戻り値にエラーが含まれていません。")
				}
			} else {
				if err != nil {
					t.Error("関数の戻り値にエラーが含まれています。")
				}
			}
		})
	}
}
