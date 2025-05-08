package utils

import (
	"fmt"
	"strings"
)

var (
	infoMessage = map[string]string{
		"INF001": "プロファイル情報が取得できないため処理を終了します。\n",
		"INF002": "クラスターが選択されていないため処理を終了します。\n",
		"INF003": "クラスターに紐づくサービスが存在しないため処理を終了します。\n",
		"INF004": "サービスが選択されていないため処理を終了します。\n",
		"INF005": "サービスに紐づくタスクが存在しないため処理を終了します。\n",
		"INF006": "タスクが選択されていないため処理を終了します。\n",
		"INF007": "タスクに紐づくコンテナが存在しないため処理を終了します。\n",
		"INF008": "コンテナが選択されていないため処理を終了します。\n",
		"INF009": "execute command が有効ではないタスクのため終了します。\n",
	}
	errorMessage = map[string]string{
		"ERR001": "session-manager-plugin がインストールされていません、以下を確認しインストールください。\nhttps://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html\n",
		"ERR002": "aws プロファイルの取得に失敗しました。\n",
		"ERR003": "該当のプロファイルに紐づく ECS クラスターが存在しないか、クラスター情報の取得に失敗しました。\n",
		"ERR004": "該当のクラスターに紐づくサービスの取得に失敗しました。\n",
		"ERR005": "該当のサービスに紐づくタスク存在しないか、タスクの取得に失敗しました。\n",
		"ERR006": "タスクに紐づくコンテナの取得に失敗しました。\n",
		"ERR999": "予期せぬエラーが発生しました。\n",
	}
	color = map[string]string{
		"default": "\x1b[30;0m",
		"red":     "\x1b[31;1m",
		"green":   "\x1b[32;5m",
	}
)

func findMessage(label string) string {
	if strings.Contains(label, "INF") {
		return color["green"] + infoMessage[label] + color["default"]
	} else {
		return color["red"] + errorMessage[label] + color["default"]
	}
}

func PrintMessage(label string) {
	fmt.Printf("%s", findMessage(label))
}
