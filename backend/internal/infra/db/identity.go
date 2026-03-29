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
)

type identityDao struct {
	PK           string `dynamodbav:"PK"`
	SK           string `dynamodbav:"SK"`
	IdentityID   string `dynamodbav:"identityId"`
	Email        string `dynamodbav:"email"`
	PasswordHash string `dynamodbav:"passwordHash"`
	GSIEmailPK   string `dynamodbav:"GSI_EMAIL_PK"`
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

// --- IdentityRepo ---

type IdentityDynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

var _ domain.IdentityRepo = (*IdentityDynamoRepo)(nil)

func NewIdentityDynamoRepo(client *dynamodb.Client, tableName string) domain.IdentityRepo {
	return &IdentityDynamoRepo{client: client, tableName: tableName}
}

func (r *IdentityDynamoRepo) FindByID(ctx context.Context, id string) (*domain.Identity, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      aws.String(r.tableName),
		ConsistentRead: aws.Bool(true),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "IDENTITY#" + id},
			"SK": &types.AttributeValueMemberS{Value: "IDENTITY"},
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

func (r *IdentityDynamoRepo) FindByEmail(ctx context.Context, email string) (*domain.Identity, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("GSI_EMAIL"),
		KeyConditionExpression: aws.String("GSI_EMAIL_PK = :email AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
			":sk":    &types.AttributeValueMemberS{Value: "IDENTITY"},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Items) == 0 {
		return nil, nil
	}
	var dao identityDao
	if err := attributevalue.UnmarshalMap(out.Items[0], &dao); err != nil {
		return nil, err
	}
	return toIdentityDomain(dao)
}

func (r *IdentityDynamoRepo) Save(ctx context.Context, identity *domain.Identity) error {
	dao := identityDao{
		PK:           "IDENTITY#" + identity.ID(),
		SK:           "IDENTITY",
		IdentityID:   identity.ID(),
		Email:        identity.Email().Value(),
		PasswordHash: string(identity.PasswordHash()),
		GSIEmailPK:   identity.Email().Value(),
		CreatedAt:    identity.CreatedAt().Format(time.RFC3339Nano),
		UpdatedAt:    identity.UpdatedAt().Format(time.RFC3339Nano),
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

func toIdentityDomain(dao identityDao) (*domain.Identity, error) {
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
	return domain.ReconstructIdentity(domain.ReconstructIdentityInput{
		ID:           id,
		Email:        dao.Email,
		PasswordHash: []byte(dao.PasswordHash),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
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
			"PK": &types.AttributeValueMemberS{Value: "SESSION#" + identityID},
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
			":pk":     &types.AttributeValueMemberS{Value: "SESSION#" + identityID},
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
		PK:        "SESSION#" + session.IdentityID(),
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
			":pk":     &types.AttributeValueMemberS{Value: "SESSION#" + identityID},
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
		end := i + transactWriteMaxItems
		if end > len(writes) {
			end = len(writes)
		}
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
