package db

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

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

// cursorKey はページネーション用の ExclusiveStartKey に含まれる属性。
// 本テーブルのキー属性はすべて文字列型。
type cursorKey struct {
	PK     string `json:"PK"`
	SK     string `json:"SK"`
	GSI1PK string `json:"GSI1PK,omitempty"`
	GSI1SK string `json:"GSI1SK,omitempty"`
	GSI2PK string `json:"GSI2PK,omitempty"`
	GSI2SK string `json:"GSI2SK,omitempty"`
	GSI3PK string `json:"GSI3PK,omitempty"`
	GSI3SK string `json:"GSI3SK,omitempty"`
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

// -----------------------------------------------------------------------
// FindByExternalID
// -----------------------------------------------------------------------

func (r *ArticleDynamoRepo) FindByExternalID(ctx context.Context, externalID string) (*domain.Article, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: articlePKPrefix + externalID},
			"SK": &types.AttributeValueMemberS{Value: articleSK},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Item) == 0 {
		return nil, nil
	}
	var dao articleDao
	if err := attributevalue.UnmarshalMap(out.Item, &dao); err != nil {
		return nil, err
	}
	return daoToArticle(dao)
}

// -----------------------------------------------------------------------
// FindAll
// -----------------------------------------------------------------------

func (r *ArticleDynamoRepo) FindAll(ctx context.Context, criteria domain.SearchCriteria) (domain.FindAllResult, error) {
	exprAttrNames := map[string]string{}
	exprAttrValues := map[string]types.AttributeValue{}

	var indexName string
	if criteria.Platform != nil {
		indexName = "GSI3"
		exprAttrNames["#pk"] = "GSI3PK"
		exprAttrValues[":pk"] = &types.AttributeValueMemberS{Value: "PLATFORM#" + *criteria.Platform}
	} else if criteria.Year != nil {
		indexName = "GSI2"
		exprAttrNames["#pk"] = "GSI2PK"
		exprAttrValues[":pk"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("YEAR#%d", *criteria.Year)}
	} else {
		indexName = "GSI1"
		exprAttrNames["#pk"] = "GSI1PK"
		exprAttrValues[":pk"] = &types.AttributeValueMemberS{Value: "ARTICLES"}
	}

	var filterParts []string

	if criteria.ActiveOnly {
		filterParts = append(filterParts, "#isActive = :isActive")
		exprAttrNames["#isActive"] = "isActive"
		exprAttrValues[":isActive"] = &types.AttributeValueMemberBOOL{Value: true}
	}
	for i, tag := range criteria.Tags {
		if i == 0 {
			exprAttrNames["#tags"] = "tags"
		}
		valKey := fmt.Sprintf(":tag%d", i)
		filterParts = append(filterParts, fmt.Sprintf("contains(#tags, %s)", valKey))
		exprAttrValues[valKey] = &types.AttributeValueMemberS{Value: tag}
	}
	for i, token := range criteria.Tokens {
		if i == 0 {
			exprAttrNames["#tokens"] = "tokens"
		}
		valKey := fmt.Sprintf(":tok%d", i)
		filterParts = append(filterParts, fmt.Sprintf("contains(#tokens, %s)", valKey))
		exprAttrValues[valKey] = &types.AttributeValueMemberS{Value: token}
	}

	var filterExpr *string
	if len(filterParts) > 0 {
		s := strings.Join(filterParts, " AND ")
		filterExpr = &s
	}

	var exclusiveStartKey map[string]types.AttributeValue
	if criteria.Cursor != nil {
		var err error
		exclusiveStartKey, err = decodeCursor(*criteria.Cursor)
		if err != nil {
			return domain.FindAllResult{}, fmt.Errorf("decode cursor: %w", err)
		}
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    aws.String("#pk = :pk"),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		FilterExpression:          filterExpr,
		ScanIndexForward:          aws.Bool(false),
		Limit:                     aws.Int32(int32(criteria.Limit)),
		ExclusiveStartKey:         exclusiveStartKey,
	}

	out, err := r.client.Query(ctx, input)
	if err != nil {
		return domain.FindAllResult{}, err
	}

	articles := make([]domain.Article, 0, len(out.Items))
	for _, item := range out.Items {
		var dao articleDao
		if err := attributevalue.UnmarshalMap(item, &dao); err != nil {
			return domain.FindAllResult{}, err
		}
		a, err := daoToArticle(dao)
		if err != nil {
			return domain.FindAllResult{}, err
		}
		articles = append(articles, *a)
	}

	result := domain.FindAllResult{Articles: articles}
	if len(out.LastEvaluatedKey) > 0 {
		cursor, err := encodeCursor(out.LastEvaluatedKey)
		if err != nil {
			return domain.FindAllResult{}, err
		}
		result.NextCursor = cursor
	}

	return result, nil
}

// -----------------------------------------------------------------------
// FindByPlatform
// -----------------------------------------------------------------------

func (r *ArticleDynamoRepo) FindByPlatform(ctx context.Context, platform string) ([]domain.Article, error) {
	var items []map[string]types.AttributeValue
	var exclusiveStartKey map[string]types.AttributeValue

	for {
		out, err := r.client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			IndexName:              aws.String("GSI3"),
			KeyConditionExpression: aws.String("#pk = :pk"),
			ExpressionAttributeNames: map[string]string{
				"#pk": "GSI3PK",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "PLATFORM#" + platform},
			},
			ExclusiveStartKey: exclusiveStartKey,
		})
		if err != nil {
			return nil, err
		}
		items = append(items, out.Items...)
		if len(out.LastEvaluatedKey) == 0 {
			break
		}
		exclusiveStartKey = out.LastEvaluatedKey
	}

	articles := make([]domain.Article, 0, len(items))
	for _, item := range items {
		var dao articleDao
		if err := attributevalue.UnmarshalMap(item, &dao); err != nil {
			return nil, err
		}
		a, err := daoToArticle(dao)
		if err != nil {
			return nil, err
		}
		articles = append(articles, *a)
	}
	return articles, nil
}

// -----------------------------------------------------------------------
// Save
// -----------------------------------------------------------------------

func (r *ArticleDynamoRepo) Save(ctx context.Context, article *domain.Article) error {
	dao := articleToDAO(article)
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

// -----------------------------------------------------------------------
// AllTags
// -----------------------------------------------------------------------

func (r *ArticleDynamoRepo) AllTags(ctx context.Context) ([]domain.TagCount, error) {
	counts := map[string]int{}
	var exclusiveStartKey map[string]types.AttributeValue

	for {
		out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:        aws.String(r.tableName),
			FilterExpression: aws.String("#isActive = :isActive"),
			ExpressionAttributeNames: map[string]string{
				"#isActive": "isActive",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":isActive": &types.AttributeValueMemberBOOL{Value: true},
			},
			ExclusiveStartKey: exclusiveStartKey,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range out.Items {
			var dao articleDao
			if err := attributevalue.UnmarshalMap(item, &dao); err != nil {
				return nil, err
			}
			for _, tag := range dao.Tags {
				counts[tag]++
			}
		}
		if len(out.LastEvaluatedKey) == 0 {
			break
		}
		exclusiveStartKey = out.LastEvaluatedKey
	}

	result := make([]domain.TagCount, 0, len(counts))
	for name, count := range counts {
		result = append(result, domain.TagCount{Name: name, Count: count})
	}
	return result, nil
}

// -----------------------------------------------------------------------
// AllTokens
// -----------------------------------------------------------------------

func (r *ArticleDynamoRepo) AllTokens(ctx context.Context) ([]domain.TokenCount, error) {
	counts := map[string]int{}
	var exclusiveStartKey map[string]types.AttributeValue

	for {
		out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:        aws.String(r.tableName),
			FilterExpression: aws.String("#isActive = :isActive"),
			ExpressionAttributeNames: map[string]string{
				"#isActive": "isActive",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":isActive": &types.AttributeValueMemberBOOL{Value: true},
			},
			ExclusiveStartKey: exclusiveStartKey,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range out.Items {
			var dao articleDao
			if err := attributevalue.UnmarshalMap(item, &dao); err != nil {
				return nil, err
			}
			for _, token := range dao.Tokens {
				counts[token]++
			}
		}
		if len(out.LastEvaluatedKey) == 0 {
			break
		}
		exclusiveStartKey = out.LastEvaluatedKey
	}

	result := make([]domain.TokenCount, 0, len(counts))
	for value, count := range counts {
		result = append(result, domain.TokenCount{Value: value, Count: count})
	}
	return result, nil
}

// -----------------------------------------------------------------------
// 変換ヘルパー
// -----------------------------------------------------------------------

func daoToArticle(dao articleDao) (*domain.Article, error) {
	var publishedAt time.Time
	if dao.PublishedAt != "" {
		t, err := time.Parse(time.RFC3339, dao.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("parse publishedAt: %w", err)
		}
		publishedAt = t
	}

	var articleUpdatedAt time.Time
	if dao.ArticleUpdatedAt != "" {
		t, err := time.Parse(time.RFC3339, dao.ArticleUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("parse articleUpdatedAt: %w", err)
		}
		articleUpdatedAt = t
	}

	createdAt, err := time.Parse(time.RFC3339, dao.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse createdAt: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, dao.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse updatedAt: %w", err)
	}

	return domain.Reconstruct(domain.ReconstructArticleInput{
		ID:               dao.ExternalID,
		Title:            dao.Title,
		URL:              dao.URL,
		Platform:         dao.Platform,
		Tags:             dao.Tags,
		Tokens:           dao.Tokens,
		PublishedAt:      publishedAt,
		ArticleUpdatedAt: articleUpdatedAt,
		IsActive:         dao.IsActive,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	})
}

func articleToDAO(a *domain.Article) articleDao {
	now := time.Now().UTC().Format(time.RFC3339)
	dao := articleDao{
		PK:         articlePKPrefix + a.ID(),
		SK:         articleSK,
		ExternalID: a.ID(),
		Title:      a.Title(),
		URL:        a.URL(),
		Platform:   a.Platform(),
		Tags:       a.Tags(),
		Tokens:     a.Tokens(),
		IsActive:   a.IsActive(),
		CreatedAt:  a.CreatedAt().UTC().Format(time.RFC3339),
		UpdatedAt:  now,
	}
	if !a.PublishedAt().IsZero() {
		dao.PublishedAt = a.PublishedAt().UTC().Format(time.RFC3339)
		dao.GSI1PK = "ARTICLES"
		dao.GSI1SK = dao.PublishedAt
		dao.GSI3PK = "PLATFORM#" + a.Platform()
		dao.GSI3SK = dao.PublishedAt
		dao.Year = a.PublishedAt().Year()
		dao.GSI2PK = fmt.Sprintf("YEAR#%d", dao.Year)
		dao.GSI2SK = dao.PublishedAt
	}
	if !a.ArticleUpdatedAt().IsZero() {
		dao.ArticleUpdatedAt = a.ArticleUpdatedAt().UTC().Format(time.RFC3339)
	}
	return dao
}

// -----------------------------------------------------------------------
// カーソルエンコード
// -----------------------------------------------------------------------

func encodeCursor(key map[string]types.AttributeValue) (*string, error) {
	ck := cursorKey{}
	if v, ok := key["PK"].(*types.AttributeValueMemberS); ok {
		ck.PK = v.Value
	}
	if v, ok := key["SK"].(*types.AttributeValueMemberS); ok {
		ck.SK = v.Value
	}
	if v, ok := key["GSI1PK"].(*types.AttributeValueMemberS); ok {
		ck.GSI1PK = v.Value
	}
	if v, ok := key["GSI1SK"].(*types.AttributeValueMemberS); ok {
		ck.GSI1SK = v.Value
	}
	if v, ok := key["GSI2PK"].(*types.AttributeValueMemberS); ok {
		ck.GSI2PK = v.Value
	}
	if v, ok := key["GSI2SK"].(*types.AttributeValueMemberS); ok {
		ck.GSI2SK = v.Value
	}
	if v, ok := key["GSI3PK"].(*types.AttributeValueMemberS); ok {
		ck.GSI3PK = v.Value
	}
	if v, ok := key["GSI3SK"].(*types.AttributeValueMemberS); ok {
		ck.GSI3SK = v.Value
	}

	b, err := json.Marshal(ck)
	if err != nil {
		return nil, err
	}
	s := base64.StdEncoding.EncodeToString(b)
	return &s, nil
}

func decodeCursor(cursor string) (map[string]types.AttributeValue, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	var ck cursorKey
	if err := json.Unmarshal(b, &ck); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	key := map[string]types.AttributeValue{}
	if ck.PK != "" {
		key["PK"] = &types.AttributeValueMemberS{Value: ck.PK}
	}
	if ck.SK != "" {
		key["SK"] = &types.AttributeValueMemberS{Value: ck.SK}
	}
	if ck.GSI1PK != "" {
		key["GSI1PK"] = &types.AttributeValueMemberS{Value: ck.GSI1PK}
	}
	if ck.GSI1SK != "" {
		key["GSI1SK"] = &types.AttributeValueMemberS{Value: ck.GSI1SK}
	}
	if ck.GSI2PK != "" {
		key["GSI2PK"] = &types.AttributeValueMemberS{Value: ck.GSI2PK}
	}
	if ck.GSI2SK != "" {
		key["GSI2SK"] = &types.AttributeValueMemberS{Value: ck.GSI2SK}
	}
	if ck.GSI3PK != "" {
		key["GSI3PK"] = &types.AttributeValueMemberS{Value: ck.GSI3PK}
	}
	if ck.GSI3SK != "" {
		key["GSI3SK"] = &types.AttributeValueMemberS{Value: ck.GSI3SK}
	}
	return key, nil
}
