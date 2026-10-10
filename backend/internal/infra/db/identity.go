package db

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	domain "github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/pkg/errs"
)

const (
	identityKeyPrefix = "IDENTITY"
	sessionKeyPrefix  = "SESSION"
)

type identityDao struct {
	PK           string `dynamodbav:"PK"`
	SK           string `dynamodbav:"SK"`
	IdentityID   string `dynamodbav:"identityId"`
	GithubUserID string `dynamodbav:"githubUserId"`
	CreatedAt    string `dynamodbav:"createdAt"`
	UpdatedAt    string `dynamodbav:"updatedAt"`
}

type sessionDao struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	UserID    string `dynamodbav:"userId"`
	TokenHash string `dynamodbav:"tokenHash"`
	Status    string `dynamodbav:"status"`
	IssuedAt  string `dynamodbav:"issuedAt"`
	ExpiresAt string `dynamodbav:"expiresAt"`
	TTL       int64  `dynamodbav:"ttl"`
}

type githubProviderDao struct {
	PK         string `dynamodbav:"PK"`
	SK         string `dynamodbav:"SK"`
	IdentityID string `dynamodbav:"identityID"`
}

// --- IdentityRepo ---

type IdentityDynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

// FindByEmail implements [identity.IdentityRepo].
func (r *IdentityDynamoRepo) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	panic("unimplemented")
}

var _ domain.IdentityRepo = (*IdentityDynamoRepo)(nil)

func NewIdentityDynamoRepo(client *dynamodb.Client, tableName string) domain.IdentityRepo {
	return &IdentityDynamoRepo{client: client, tableName: tableName}
}

// FindByGithubID implements [identity.IdentityRepo].
func (r *IdentityDynamoRepo) FindByGithubID(ctx context.Context, githubID string) (*domain.Account, error) {
	out, err := r.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: "GITHUB" + "#" + githubID},
				"SK": &types.AttributeValueMemberS{Value: "LOOKUP"},
			},
			TableName:      aws.String(r.tableName),
			ConsistentRead: aws.Bool(true),
		},
	)
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, errs.ErrNotFound
	}
	var dao githubProviderDao
	if err := attributevalue.UnmarshalMap(out.Item, &dao); err != nil {
		return nil, err
	}
	outIdentity, err := r.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: identityKeyPrefix + "#" + dao.IdentityID},
				"SK": &types.AttributeValueMemberS{Value: identityKeyPrefix},
			},
			TableName:      aws.String(r.tableName),
			ConsistentRead: aws.Bool(true),
		},
	)
	if err != nil {
		return nil, err
	}
	var daoIdentity identityDao
	if err := attributevalue.UnmarshalMap(outIdentity.Item, &daoIdentity); err != nil {
		return nil, err
	}
	return toIdentityDomain(daoIdentity)
}

func (r *IdentityDynamoRepo) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      aws.String(r.tableName),
		ConsistentRead: aws.Bool(true),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: identityKeyPrefix + "#" + id},
			"SK": &types.AttributeValueMemberS{Value: identityKeyPrefix},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	var dao identityDao
	if err := attributevalue.UnmarshalMap(out.Item, &dao); err != nil {
		return nil, err
	}
	return toIdentityDomain(dao)
}

func (r *IdentityDynamoRepo) Save(ctx context.Context, identity *domain.Account) error {
	dao := identityDao{
		PK:           identityKeyPrefix + "#" + identity.ID(),
		SK:           identityKeyPrefix,
		IdentityID:   identity.ID(),
		CreatedAt:    identity.CreatedAt().Format(time.RFC3339Nano),
		UpdatedAt:    identity.UpdatedAt().Format(time.RFC3339Nano),
		GithubUserID: identity.GithubID(),
	}
	item, err := attributevalue.MarshalMap(dao)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{TableName: aws.String(r.tableName), Item: item},
	)
	if err != nil {
		return err
	}
	if dao.GithubUserID != "" {
		githubDao := githubProviderDao{
			PK:         "GITHUB#" + dao.GithubUserID,
			SK:         "LOOKUP",
			IdentityID: dao.IdentityID,
		}
		githubItem, err := attributevalue.MarshalMap(githubDao)
		if err != nil {
			return err
		}
		_, err = r.client.PutItem(
			ctx,
			&dynamodb.PutItemInput{TableName: aws.String(r.tableName), Item: githubItem},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func toIdentityDomain(dao identityDao) (*domain.Account, error) {
	id, err := uuid.Parse(dao.IdentityID)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339Nano, dao.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("createdAt parse error: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, dao.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updatedAt parse error: %w", err)
	}
	return domain.ReconstructAccount(domain.ReconstructIdentityInput{
		InputID:        id,
		InputCreatedAt: createdAt,
		InputUpdatedAt: updatedAt,
		InputGithubID:  dao.GithubUserID,
	})
}

// --- SessionRepo ---

type SessionDynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

var _ domain.SessionRepo = (*SessionDynamoRepo)(nil)

func NewSessionDynamoRepo(client *dynamodb.Client, tableName string) domain.SessionRepo {
	return &SessionDynamoRepo{client: client, tableName: tableName}
}

func (r *SessionDynamoRepo) FindByIdentityIdAndTokenHash(ctx context.Context, identityID, tokenHash string) (*domain.Session, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      aws.String(r.tableName),
		ConsistentRead: aws.Bool(true),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: sessionKeyPrefix + "#" + identityID},
			"SK": &types.AttributeValueMemberS{Value: "RT#" + tokenHash},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	var dao sessionDao
	if err := attributevalue.UnmarshalMap(out.Item, &dao); err != nil {
		return nil, err
	}
	return toSessionDomain(dao)
}

func (r *SessionDynamoRepo) FindActiveByIdentity(ctx context.Context, identityID string) ([]*domain.Session, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		FilterExpression:       aws.String("#st = :active"),
		ExpressionAttributeNames: map[string]string{
			"#st": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: sessionKeyPrefix + "#" + identityID},
			":prefix": &types.AttributeValueMemberS{Value: "RT#"},
			":active": &types.AttributeValueMemberS{Value: "active"},
		},
	})
	if err != nil {
		return nil, err
	}
	sessions := make([]*domain.Session, 0, len(out.Items))
	for _, item := range out.Items {
		var dao sessionDao
		if err := attributevalue.UnmarshalMap(item, &dao); err != nil {
			return nil, err
		}
		s, err := toSessionDomain(dao)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *SessionDynamoRepo) Save(ctx context.Context, session *domain.Session) error {
	dao := sessionDao{
		PK:        sessionKeyPrefix + "#" + session.IdentityID(),
		SK:        "RT#" + session.TokenHash(),
		UserID:    session.IdentityID(),
		TokenHash: session.TokenHash(),
		Status:    session.Status(),
		IssuedAt:  session.IssuedAt().Format(time.RFC3339Nano),
		ExpiresAt: session.ExpiresAt().Format(time.RFC3339Nano),
		TTL:       session.ExpiresAt().Unix(),
	}
	item, err := attributevalue.MarshalMap(dao)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

const transactWriteMaxItems = 25

func (r *SessionDynamoRepo) RevokeAll(ctx context.Context, identityID string) error {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		ConsistentRead:         aws.Bool(true),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: sessionKeyPrefix + "#" + identityID},
			":prefix": &types.AttributeValueMemberS{Value: "RT#"},
		},
	})
	if err != nil {
		return err
	}

	var writes []types.TransactWriteItem
	for _, item := range out.Items {
		pkAttr, ok := item["PK"].(*types.AttributeValueMemberS)
		if !ok {
			return fmt.Errorf("invalid PK attribute type")
		}
		skAttr, ok := item["SK"].(*types.AttributeValueMemberS)
		if !ok {
			return fmt.Errorf("invalid SK attribute type")
		}
		writes = append(writes, types.TransactWriteItem{
			Update: &types.Update{
				TableName: aws.String(r.tableName),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: pkAttr.Value},
					"SK": &types.AttributeValueMemberS{Value: skAttr.Value},
				},
				UpdateExpression: aws.String("SET #st = :revoked"),
				ExpressionAttributeNames: map[string]string{
					"#st": "status",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":revoked": &types.AttributeValueMemberS{Value: "revoked"},
				},
			},
		})
	}

	for i := 0; i < len(writes); i += transactWriteMaxItems {
		end := min(i+transactWriteMaxItems, len(writes))
		_, err := r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: writes[i:end],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func toSessionDomain(dao sessionDao) (*domain.Session, error) {
	issuedAt, err := time.Parse(time.RFC3339Nano, dao.IssuedAt)
	if err != nil {
		return nil, fmt.Errorf("issuedAt parse error: %w", err)
	}
	expiresAt, err := time.Parse(time.RFC3339Nano, dao.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("expiresAt parse error: %w", err)
	}
	return domain.ReconstructSession(domain.ReconstructSessionInput{
		IdentityID: dao.UserID,
		TokenHash:  dao.TokenHash,
		Status:     dao.Status,
		IssuedAt:   issuedAt,
		ExpiresAt:  expiresAt,
	})
}
