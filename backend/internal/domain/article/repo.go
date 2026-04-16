package article

import "context"

type Repo interface {
	FindByExternalID(
		ctx context.Context,
		externalID string,
	) (*Article, error)
	FindAll(
		ctx context.Context,
		criteria SearchCriteria,
	) (FindAllResult, error)
	FindByPlatform(
		ctx context.Context,
		platform string,
	) ([]Article, error)
	Save(ctx context.Context, article Article) error
	Exists(ctx context.Context, externalID string) (bool, error)
	AllTags(ctx context.Context) ([]TagCount, error)
	AllTokens(ctx context.Context) ([]TokenCount, error)
}
