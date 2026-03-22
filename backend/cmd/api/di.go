package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	domain "github.com/umekikazuya/me/internal/domain/me"
	inframe "github.com/umekikazuya/me/internal/infra/me"
)

// setupRepo はリポジトリの依存関係を初期化する
func setupRepo(ctx context.Context) (domain.Repo, error) {
	// AWS 設定の読み込み
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// DynamoDB クライアントの生成
	client := dynamodb.NewFromConfig(cfg)

	// テーブル名の取得
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "me"
	}

	return inframe.NewDynamoRepo(client, tableName), nil
}
