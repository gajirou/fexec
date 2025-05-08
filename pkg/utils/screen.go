package utils

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var (
	labelMessage = map[string]string{
		"cluster":   "対象のクラスター名を選択してください：",
		"service":   "対象のサービス名を選択してください：",
		"task":      "対象のタスク ID を選択してください：",
		"container": "対象のコンテナを選択してください：",
	}
)

var answers struct {
	Askone string `survey:"askone"`
}

func ScreenDraw(options []string, label string) (string, error) {
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
		if err == terminal.InterruptErr {
			return "", nil
		}
		return "", err
	}
	return answers.Askone, nil
}
