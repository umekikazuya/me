package article

import "context"

type Repo interface {
	FindByExternalID(ctx context.Context, externalID string) (*Article, error)
	FindAll(
		ctx context.Context,
		criteria SearchCriteria,
	) (
		[]Article,
		*string,
		error,
	)
	FindByPlatform(
		ctx context.Context,
		plaplatform string,
	) (
		[]Article,
		int,
		error,
	)
	Save(ctx context.Context, article Article) error
	Exists(ctx context.Context, externalID string) (bool, error)
	AllTags(ctx context.Context) ([]TagCount, error)
	AllTokens(ctx context.Context) ([]TokenCount, error)
}
