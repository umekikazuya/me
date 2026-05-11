package di

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	apparticle "github.com/umekikazuya/me/internal/app/article"
	"github.com/umekikazuya/me/internal/app/eventhandler"
	appidentity "github.com/umekikazuya/me/internal/app/identity"
	appme "github.com/umekikazuya/me/internal/app/me"
	"github.com/umekikazuya/me/internal/domain/article"
	"github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/internal/domain/me"
	handlerarticle "github.com/umekikazuya/me/internal/handler/article"
	handleridentity "github.com/umekikazuya/me/internal/handler/identity"
	handlerme "github.com/umekikazuya/me/internal/handler/me"
	"github.com/umekikazuya/me/internal/infra/db"
	infraevent "github.com/umekikazuya/me/internal/infra/event"
	"github.com/umekikazuya/me/internal/infra/fetcher"
	"github.com/umekikazuya/me/internal/infra/token"
	"github.com/umekikazuya/me/internal/infra/tokenizer"
)

func setupRepo(ctx context.Context) (me.Repo, identity.IdentityRepo, identity.SessionRepo, article.Repo, error) {
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
	meRepo := db.NewMeDynamoRepo(client, tableName)
	identityRepo := db.NewIdentityDynamoRepo(client, tableName)
	sessionRepo := db.NewSessionDynamoRepo(client, tableName)

	return meRepo, identityRepo, sessionRepo, articleRepo, nil
}

func NewHandlers(ctx context.Context) (*Handlers, error) {
	// Repo
	meRepo, identityRepo, sessionRepo, articleRepo, err := setupRepo(ctx)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"インフラの初期化に失敗しました",
			"error",
			err,
		)
		return nil, err
	}
	// 環境変数
	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		slog.ErrorContext(ctx, "JWT_SECRET が未設定です")
		return nil, errors.New("JWT_SECRET is not set")
	}

	// 各アダプタ
	tokenSrv := token.NewJWTTokenService(
		jwtSecret,
		15*time.Minute,
	)
	articleFetcher := fetcher.NewDefaultDispatcher(
		os.Getenv("QIITA_TOKEN"),
		os.Getenv("ZENN_USERNAME"),
	)
	articleTokenizer, err := tokenizer.NewKagomeTokenizer()
	if err != nil {
		slog.ErrorContext(ctx, "トークナイザーの初期化に失敗しました", "error", err)
		return nil, err
	}

	articleInteractor := apparticle.NewInteractor(
		articleRepo,
		articleFetcher,
		articleTokenizer,
	)
	meInteractor := appme.NewInteractor(meRepo)

	// ディスパッチャー
	dispatcher := infraevent.NewSyncEventDispatcher()
	dispatcher.Register(eventhandler.NewIdentityRegisteredHandler(meInteractor))

	// ユースケース
	identityInteractor := appidentity.NewInteractor(identityRepo, sessionRepo, tokenSrv, dispatcher)
	return &Handlers{
		Me:       *handlerme.NewHandler(meInteractor),
		Article:  *handlerarticle.NewHandler(articleInteractor),
		Identity: *handleridentity.NewHandler(identityInteractor, tokenSrv),
	}, nil
}
