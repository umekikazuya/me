package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	domain "github.com/umekikazuya/me/internal/domain/article"
)

const (
	articlePKPrefix = "ARTICLE#"
	articleSK       = "ARTICLE"
)

type articleDao struct {
	PK               string   `dynamodbav:"PK"`
	SK               string   `dynamodbav:"SK"`
	GSI1PK           string   `dynamodbav:"GSI1PK,omitempty"`
	GSI1SK           string   `dynamodbav:"GSI1SK,omitempty"`
	GSI2PK           string   `dynamodbav:"GSI2PK,omitempty"`
	GSI2SK           string   `dynamodbav:"GSI2SK,omitempty"`
	GSI3PK           string   `dynamodbav:"GSI3PK,omitempty"`
	GSI3SK           string   `dynamodbav:"GSI3SK,omitempty"`
	ExternalID       string   `dynamodbav:"externalId"`
	Title            string   `dynamodbav:"title"`
	URL              string   `dynamodbav:"url"`
	Platform         string   `dynamodbav:"platform"`
	PublishedAt      string   `dynamodbav:"publishedAt,omitempty"`
	ArticleUpdatedAt string   `dynamodbav:"articleUpdatedAt,omitempty"`
	Tags             []string `dynamodbav:"tags,omitempty"`
	Tokens           []string `dynamodbav:"tokens,omitempty"`
	IsActive         bool     `dynamodbav:"isActive"`
	Year             int      `dynamodbav:"year,omitempty"`
	CreatedAt        string   `dynamodbav:"createdAt"`
	UpdatedAt        string   `dynamodbav:"updatedAt"`
}

type ArticleDynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

var _ domain.Repo = (*ArticleDynamoRepo)(nil)

func NewArticleDynamoRepo(client *dynamodb.Client, tableName string) domain.Repo {
	return &ArticleDynamoRepo{
		client:    client,
		tableName: tableName,
	}
}

func (r *ArticleDynamoRepo) FindByExternalID(ctx context.Context, externalID string) (*domain.Article, error) {
	return nil, errors.New("not implemented")
}

func (r *ArticleDynamoRepo) FindAll(ctx context.Context, criteria domain.SearchCriteria) (domain.FindAllResult, error) {
	return domain.FindAllResult{}, errors.New("not implemented")
}

func (r *ArticleDynamoRepo) FindByPlatform(ctx context.Context, platform string) ([]domain.Article, error) {
	return nil, errors.New("not implemented")
}

func (r *ArticleDynamoRepo) Save(ctx context.Context, article *domain.Article) error {
	return errors.New("not implemented")
}

func (r *ArticleDynamoRepo) AllTags(ctx context.Context) ([]domain.TagCount, error) {
	return nil, errors.New("not implemented")
}

func (r *ArticleDynamoRepo) AllTokens(ctx context.Context) ([]domain.TokenCount, error) {
	return nil, errors.New("not implemented")
}
