# Nunu — Goアプリケーションを構築するためのCLIツール

Nunuは、Goアプリケーションを構築するためのスキャフォールディングツールです。その名前は、リーグ・オブ・レジェンドのゲームキャラクターから来ており、イエティの肩に乗る小さな男の子を意味します。Nunuのように、このプロジェクトも巨人の肩の上に立っており、Goエコシステムからの人気のあるライブラリの組み合わせに基づいて構築されています。この組み合わせにより、効率的で信頼性の高いアプリケーションを迅速に構築できます。

🚀ヒント: このプロジェクトは非常に完成度が高いため、更新は頻繁ではありませんが、ぜひご利用ください。

- [英語版](https://github.com/go-nunu/nunu/blob/main/README.md)
- [簡体中文版](https://github.com/go-nunu/nunu/blob/main/README_zh.md)
- [ポルトガル語版](https://github.com/go-nunu/nunu/blob/main/README_pt.md)

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)

## ドキュメント
* [ユーザーガイド](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [アーキテクチャ](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [入門チュートリアル](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)
* [ユニットテスト](https://github.com/go-nunu/nunu/blob/main/docs/en/unit_testing.md)

## 特徴
- **Gin**: https://github.com/gin-gonic/gin
- **Gorm**: https://github.com/go-gorm/gorm
- **Wire**: https://github.com/google/wire
- **Viper**: https://github.com/spf13/viper
- **Zap**: https://github.com/uber-go/zap
- **Golang-jwt**: https://github.com/golang-jwt/jwt
- **Go-redis**: https://github.com/go-redis/redis
- **Testify**: https://github.com/stretchr/testify
- **Sonyflake**: https://github.com/sony/sonyflake
- **Gocron**:  https://github.com/go-co-op/gocron
- **Go-sqlmock**:  https://github.com/DATA-DOG/go-sqlmock
- **Gomock**:  https://github.com/golang/mock
- **Swaggo**:  https://github.com/swaggo/swag
- **Pitaya**:  https://github.com/topfreegames/pitaya
- **Casbin**:  https://github.com/casbin/casbin

- その他多数...

## 主な特徴
* **低い学習曲線とカスタマイズ性**: Nunuは、Gopherがよく知る人気のあるライブラリをカプセル化しており、特定の要件を満たすためにアプリケーションを簡単にカスタマイズできます。
* **高性能とスケーラビリティ**: Nunuは、高性能でスケーラブルを目指しています。最新の技術とベストプラクティスを使用して、アプリケーションが高いトラフィックと大量のデータを処理できるようにしています。
* **セキュリティと信頼性**: Nunuは、安定して信頼性の高いサードパーティのライブラリを使用して、アプリケーションのセキュリティと信頼性を保証しています。
* **モジュラー性と拡張性**: Nunuは、モジュラー性と拡張性を目指して設計されています。サードパーティのライブラリを使用するか、独自のモジュールを書くことで、新しい機能と機能を簡単に追加できます。
* **完全なドキュメントとテスト**: Nunuには、包括的なドキュメントとテストがあります。広範なドキュメントと例を提供して、迅速に開始できるようにしています。また、アプリケーションが期待通りに動作することを保証するテストスイートも含まれています。

## 簡潔なレイヤードアーキテクチャ
Nunuは、クラシックなレイヤードアーキテクチャを採用しています。モジュラリティとデカップリングを実現するために、依存性注入フレームワーク`Wire`を使用しています。

![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.png)

## Nunu CLI

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/screenshot.jpg)


## ディレクトリ構造
```
.
├── api
│   └── v1
├── cmd
│   ├── migration
│   ├── server
│   │   ├── wire
│   │   │   ├── wire.go
│   │   │   └── wire_gen.go
│   │   └── main.go
│   └── task
├── config
├── deploy
├── docs
├── internal
│   ├── handler
│   ├── middleware
│   ├── model
│   ├── repository
│   ├── server
│   └── service
├── pkg
├── scripts
├── test
│   ├── mocks
│   └── server
├── web
├── Makefile
├── go.mod
└── go.sum

```

プロジェクトのアーキテクチャは、典型的なレイヤード構造に従っており、以下のモジュールで構成されています：

* `cmd`: このモジュールには、アプリケーションのエントリーポイントが含まれており、異なるコマンドに基づいて異なる操作を実行します（例：サーバーの起動、データベースマイグレーションの実行など）。各サブモジュールには、エントリーファイルとしての`main.go`ファイルと、依存性注入のための`wire.go`および`wire_gen.go`ファイルが含まれています。
* `config`: このモジュールには、アプリケーションの設定ファイルが含まれており、開発環境と本番環境など、異なる環境に対する異なる設定を提供します。
* `deploy`: このモジュールは、アプリケーションのデプロイに使用され、デプロイスクリプトと設定ファイルが含まれています。
* `internal`: このモジュールは、アプリケーションのコアモジュールであり、さまざまなビジネスロジックの実装が含まれています。

  - `handler`: このサブモジュールには、HTTPリクエストを処理するハンドラーが含まれており、リクエストを受け取り、対応するサービスを呼び出して処理します。

  - `job`: このサブモジュールには、バックグラウンドタスクのロジックが含まれています。

  - `model`: このサブモジュールには、データモデルの定義が含まれています。

  - `repository`: このサブモジュールには、データアクセス層の実装が含まれており、データベースとのやり取りを担当します。

  - `server`: このサブモジュールには、HTTPサーバーの実装が含まれています。

  - `service`: このサブモジュールには、ビジネスロジックの実装が含まれており、特定のビジネス操作を処理します。

* `pkg`: このモジュールには、いくつかの共通のユーティリティと機能が含まれています。

* `scripts`: このモジュールには、プロジェクトのビルド、テスト、デプロイメント操作用のスクリプトファイルが含まれています。

* `storage`: このモジュールは、ファイルやその他の静的リソースを保存するために使用されます。

* `test`: このモジュールには、さまざまなモジュールのユニットテストが含まれており、モジュールに基づいてサブディレクトリに整理されています。

* `web`: このモジュールには、HTML、CSS、JavaScriptなど、フロントエンド関連のファイルが含まれています。

さらに、ライセンスファイル、ビルドファイル、READMEなど、いくつかの他のファイルとディレクトリが含まれています。全体として、プロジェクトのアーキテクチャは明確であり、各モジュールの責任が明確であり、理解とメンテナンスが容易です。

## 要件
Nunuを使用するには、システムに次のソフトウェアがインストールされている必要があります：

* Go 1.19以上
* Git
* Docker（オプション）
* MySQL 5.7以上（オプション）
* Redis（オプション）

### インストール

次のコマンドでNunuをインストールできます：

```bash
go install github.com/go-nunu/nunu@latest
```

> ヒント: `go install`が成功しても`nunu`コマンドが認識されない場合は、環境変数が設定されていないためです。GOBINディレクトリを環境変数に追加してください。

### 新しいプロジェクトを作成する

次のコマンドで新しいGoプロジェクトを作成できます：

```bash
nunu new projectName
```

デフォルトでは、GitHubソースからプルしますが、中国で加速されたリポジトリを使用することもできます：

```
// 基本テンプレートを使用する
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-basic.git
// 高度なテンプレートを使用する
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-advanced.git
```

このコマンドは`projectName`という名前のディレクトリを作成し、その中にエレガントなGoプロジェクト構造を生成します。

### コンポーネントを作成する

次のコマンドを使用して、プロジェクトのハンドラー、サービス、リポジトリ、モデルなどのコンポーネントを作成できます：

```bash
nunu create handler user
nunu create service user
nunu create repository user
nunu create model user
```
または
```
nunu create all user
```

これらのコマンドは、それぞれ`UserHandler`、`UserService`、`UserRepository`、`UserModel`という名前のコンポーネントを作成し、正しいディレクトリに配置します。

### プロジェクトを実行する

次のコマンドでプロジェクトをすぐに実行できます：

```bash
nunu run
```

このコマンドは、Goプロジェクトを起動し、ファイルが更新されたときにホットリロードをサポートします。

### wire.goをコンパイルする

次のコマンドで`wire.go`をすぐにコンパイルできます：

```bash
nunu wire
```

このコマンドは、`wire.go`ファイルをコンパイルし、必要な依存関係を生成します。

## 貢献

問題を発見したり、改善提案がある場合は、遠慮なく問題を提起したり、プルリクエストを送信してください。あなたの貢献を大歓迎します！

## ライセンス

Nunuは、MITライセンスの下でリリースされています。詳細については、[LICENSE](LICENSE)ファイルを参照してください。

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=go-nunu/nunu&type=Date)](https://star-history.com/#go-nunu/nunu&Date)
