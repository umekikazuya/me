package appevent

import (
	"context"

	pkgdomain "github.com/umekikazuya/me/pkg/domain"
)

// EventHandler はドメインイベントを処理するハンドラーのインターフェース
type EventHandler interface {
	EventType() string
	Handle(ctx context.Context, event pkgdomain.DomainEvent) error
}

// EventDispatcher はイベントのディスパッチを担うアプリケーション層のポート
type EventDispatcher interface {
	Dispatch(ctx context.Context, events []pkgdomain.DomainEvent) error
}
