package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/internal/domain/me"
	"github.com/umekikazuya/me/internal/infra/db"
)

// setupRepo はリポジトリの依存関係を初期化する
func setupRepo(ctx context.Context) (me.Repo, identity.IdentityRepo, identity.SessionRepo, error) {
	// AWS 設定の読み込み
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	// DynamoDB クライアントの生成
	client := dynamodb.NewFromConfig(cfg)

	// テーブル名の取得
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "me"
	}

	return db.NewMeDynamoRepo(client, tableName), db.NewIdentityDynamoRepo(client, tableName), db.NewSessionDynamoRepo(client, tableName), nil
}
