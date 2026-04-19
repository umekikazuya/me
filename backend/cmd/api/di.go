package main

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	apparticle "github.com/umekikazuya/me/internal/app/article"
	"github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/internal/domain/me"
	"github.com/umekikazuya/me/internal/infra/db"
	"github.com/umekikazuya/me/internal/infra/fetcher"
)

// setupRepo はリポジトリの依存関係を初期化する
func setupRepo(ctx context.Context) (me.Repo, identity.IdentityRepo, identity.SessionRepo, apparticle.Interactor, error) {
	loadCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(loadCtx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// DynamoDB クライアントの生成
	client := dynamodb.NewFromConfig(cfg)

	// テーブル名の取得
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "me"
	}

	articleRepo := db.NewArticleDynamoRepo(client, tableName)
	articleFetcher := fetcher.NewDefaultDispatcher(
		os.Getenv("QIITA_TOKEN"),
		os.Getenv("ZENN_USERNAME"),
	)
	articleInteractor := apparticle.NewInteractor(articleRepo, articleFetcher)

	return db.NewMeDynamoRepo(client, tableName), db.NewIdentityDynamoRepo(client, tableName), db.NewSessionDynamoRepo(client, tableName), articleInteractor, nil
}
