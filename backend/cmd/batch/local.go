package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	apparticle "github.com/umekikazuya/me/internal/app/article"
	"github.com/umekikazuya/me/internal/infra/db"
	"github.com/umekikazuya/me/internal/infra/fetcher"
	"github.com/umekikazuya/me/internal/infra/tokenizer"
)

var targetPlatforms = []string{"qiita", "zenn"}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	ctx := context.Background()

	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "me"
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
	)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	var clientOpts []func(*dynamodb.Options)
	if endpoint != "" {
		clientOpts = append(clientOpts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}
	client := dynamodb.NewFromConfig(cfg, clientOpts...)

	articleRepo := db.NewArticleDynamoRepo(client, tableName)
	articleFetcher := fetcher.NewDefaultDispatcher(
		os.Getenv("QIITA_TOKEN"),
		os.Getenv("ZENN_USERNAME"),
	)
	articleTokenizer, err := tokenizer.NewKagomeTokenizer()
	if err != nil {
		slog.Error("failed to init tokenizer", "error", err)
		os.Exit(1)
	}
	interactor := apparticle.NewInteractor(articleRepo, articleFetcher, articleTokenizer)

	hasError := false
	for _, platform := range targetPlatforms {
		slog.InfoContext(ctx, "syncing", "platform", platform)
		result := interactor.Sync(ctx, platform)

		errMsgs := make([]string, 0, len(result.Errors))
		for _, e := range result.Errors {
			errMsgs = append(errMsgs, e.Error())
		}

		if len(result.Errors) > 0 {
			hasError = true
			slog.ErrorContext(ctx, "sync completed with errors",
				"platform", platform,
				"indexed", result.Indexed,
				"reindexed", result.Reindexed,
				"deactivated", result.Deactivated,
				"errors", errMsgs,
			)
		} else {
			slog.InfoContext(ctx, "sync completed",
				"platform", platform,
				"indexed", result.Indexed,
				"reindexed", result.Reindexed,
				"deactivated", result.Deactivated,
			)
		}
	}

	if hasError {
		os.Exit(1)
	}
}
