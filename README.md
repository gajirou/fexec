# fexec（仮）
AWS Fargate で動作しているコンテナへ接続するだけのコマンド。
## 前提条件
session-manager-plugin がインストールされている事、詳細は以下を参照ください。

https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html
## インストール
### macOS
```
brew tap gajirou/fexec
brew install fexec
```
### Linux
```
wget https://github.com/gajirou/fexec/releases/latest/download/fexec_linux_amd64.tar.gz
tar xzf fexec*.tar.gz
```
## 利用方法
環境変数 `AWS_SESSION_TOKEN` が設定されている場合はセッション情報を利用。

設定がない場合は、環境変数 `AWS_DEFAULT_PROFILE` に設定されているプロファイル名から AWS Credencial 情報を取得し、取得できない場合は default の AWS Credencial を利用する。

また、パラーメータで指定したプロファイル情報の指定も可能。

クラスター名、サービス名、タスク ARN、コンテナ名を選択 or 入力すると、execute command が有効の場合に該当コンテナに接続する。

![fexec](https://storage.googleapis.com/zenn-user-upload/3013879517cb-20220806.gif)

## パラメータ
| パラメータ | 設定値 |
| ---- | ---- |
| -p | 利用プロファイル名（初期値：default） |
## 今後やる
- リファクタリング
- windows コンテナ対応
- EC2 タイプコンテナ
- Readme をかっこよくする
- AWS の各リソース取得上限をいい感じに処理できるようにする。
