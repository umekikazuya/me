package me

import (
	"context"
	"time"

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

type linkDao struct {
	Platform string `dynamodbav:"platform"`
	URL      string `dynamodbav:"url"`
}

type certificationDao struct {
	Issuer string `dynamodbav:"issuer,omitempty"`
	Month  int    `dynamodbav:"month"`
	Name   string `dynamodbav:"name"`
	Year   int    `dynamodbav:"year"`
}

type meDao struct {
	PK             string             `dynamodbav:"PK"`
	SK             string             `dynamodbav:"SK"`
	DisplayName    string             `dynamodbav:"display"`
	DisplayNameJa  string             `dynamodbav:"displayJa,omitempty"`
	Role           string             `dynamodbav:"role,omitempty"`
	Location       string             `dynamodbav:"location,omitempty"`
	Likes          []string           `dynamodbav:"likes,omitempty"`
	Links          []linkDao          `dynamodbav:"links,omitempty"`
	Certifications []certificationDao `dynamodbav:"certifications,omitempty"`
	CreatedAt      string             `dynamodbav:"createdAt,omitempty"`
	UpdatedAt      string             `dynamodbav:"updatedAt,omitempty"`
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

	createdAt, _ := time.Parse(time.RFC3339Nano, dao.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, dao.UpdatedAt)

	links := make([]domain.Link, 0, len(dao.Links))
	for _, l := range dao.Links {
		link, err := domain.NewLink(l.Platform, l.URL)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	certifications := make([]domain.Certification, 0, len(dao.Certifications))
	for _, v := range dao.Certifications {
		certification, err := domain.NewCertification(v.Name, v.Issuer, v.Year, v.Month)
		if err != nil {
			return nil, err
		}
		certifications = append(certifications, certification)
	}

	input := domain.ReconstructInput{
		Name:           dao.DisplayName,
		Likes:          dao.Likes,
		Links:          links,
		Certifications: certifications,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
	if dao.DisplayNameJa != "" {
		input.DisplayJa = &dao.DisplayNameJa
	}
	if dao.Role != "" {
		input.Role = &dao.Role
	}
	if dao.Location != "" {
		input.Location = &dao.Location
	}

	return domain.Reconstruct(input), nil
}

func (repo *DynamoRepo) Save(ctx context.Context, me *domain.Me) error {
	links := make([]linkDao, 0, len(me.Links()))
	for _, l := range me.Links() {
		links = append(links, linkDao{Platform: l.Platform(), URL: l.URL()})
	}
	c := make([]certificationDao, 0, len(me.Certifications()))
	for _, v := range me.Certifications() {
		c = append(c, certificationDao{
			Name:   v.Name(),
			Issuer: v.Issuer(),
			Year:   v.Year(),
			Month:  v.Month(),
		})
	}

	dao := meDao{
		PK:             profilePK,
		SK:             profileSK,
		DisplayName:    me.DisplayName(),
		DisplayNameJa:  me.DisplayNameJa(),
		Role:           me.Role(),
		Location:       me.Location(),
		Likes:          me.Likes(),
		Links:          links,
		Certifications: c,
		CreatedAt:      me.CreatedAt().Format(time.RFC3339Nano),
		UpdatedAt:      me.UpdatedAt().Format(time.RFC3339Nano),
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
