package me

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	domain "github.com/umekikazuya/me/internal/domain/me"
)

const (
	profilePK = "PROFILE"
	profileSK = "PROFILE"
)

type meDao struct {
	PK            string   `dynamodbav:"PK"`
	SK            string   `dynamodbav:"SK"`
	DisplayName   string   `dynamodbav:"display"`
	DisplayNameJa string   `dynamodbav:"displayJa,omitempty"`
	Role          string   `dynamodbav:"role,omitempty"`
	Location      string   `dynamodbav:"location,omitempty"`
	Likes         []string `dynamodbav:"likes,omitempty"`
	UpdatedAt     string   `dynamodbav:"updatedAt,omitempty"`
}

type DynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

var _ domain.Repo = (*DynamoRepo)(nil)

func NewDynamoRepo(client *dynamodb.Client, tableName string) domain.Repo {
	return &DynamoRepo{
		client:    client,
		tableName: tableName,
	}
}

func (repo *DynamoRepo) Find(ctx context.Context) (*domain.Me, error) {
	out, err := repo.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: profilePK},
			"SK": &types.AttributeValueMemberS{Value: profileSK},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, nil
	}

	var dao meDao
	if err := attributevalue.UnmarshalMap(out.Item, &dao); err != nil {
		return nil, err
	}

	// ドメインエンティティへの変換 (Reconstruct)
	var opts []domain.OptFunc
	if dao.DisplayNameJa != "" {
		opts = append(opts, domain.OptDisplayNameJa(dao.DisplayNameJa))
	}
	if dao.Role != "" {
		opts = append(opts, domain.OptRole(dao.Role))
	}
	if dao.Location != "" {
		opts = append(opts, domain.OptLocation(dao.Location))
	}
	if len(dao.Likes) > 0 {
		opts = append(opts, domain.OptLikes(dao.Likes))
	}

	return domain.NewMe(dao.DisplayName, opts...)
}

func (repo *DynamoRepo) Save(ctx context.Context, me *domain.Me) error {
	dao := meDao{
		PK:            profilePK,
		SK:            profileSK,
		DisplayName:   me.DisplayName(),
		DisplayNameJa: me.DisplayNameJa(),
		Role:          me.Role(),
		Location:      me.Location(),
		Likes:         me.Likes(),
	}

	item, err := attributevalue.MarshalMap(dao)
	if err != nil {
		return err
	}

	_, err = repo.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item:      item,
	})
	return err
}

func (repo *DynamoRepo) Exists(ctx context.Context) (bool, error) {
	out, err := repo.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: profilePK},
			"SK": &types.AttributeValueMemberS{Value: profileSK},
		},
		ProjectionExpression: aws.String("PK"),
	})
	if err != nil {
		return false, err
	}

	return out.Item != nil, nil
}
