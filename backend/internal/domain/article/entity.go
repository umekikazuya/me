package article

import "time"

type Article struct {
	id               id
	title            title
	platform         platform
	publishedAt      publishedAt
	articleUpdatedAt articleUpdatedAt
	isActive         isActive
	createdAt        time.Time
	updatedAt        time.Time
}
