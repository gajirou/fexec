package fexec

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var (
	// ターミナル描画用ラベルメッセージ
	labelMessage = map[string]string{
		"cluster":   "対象のクラスター名を選択してください：",
		"service":   "対象のサービス名を選択してください：",
		"task":      "対象のタスク ID を選択してください：",
		"container": "対象のコンテナを選択してください：",
	}
)

// 回答保持用構造体
var answers struct {
	Askone string `survey:"askone"`
}

// ターミナルリスト選択
func ScreenDraw(options []string, label string) (string, error) {
	// select 形式で表示
	var qs = []*survey.Question{
		{
			Name: "askone",
			Prompt: &survey.Select{
				Message: labelMessage[label],
				Options: options,
				Default: options[0],
			},
			Validate: survey.Required,
		},
	}
	err := survey.Ask(qs, &answers)
	if err != nil {
		// SIGINT の場合は err を返さない
		if err == terminal.InterruptErr {
			return "", nil
		}
		PrintMessage("ERR999")
		return "", err
	}
	return answers.Askone, nil
}
