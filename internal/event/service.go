package event

import "context"

type Service interface {
	GetById(ctx context.Context, id string) (*Event, error)
	GetAll(ctx context.Context, page, limit int) ([]*Event, error)
	Insert(ctx context.Context, in *Event) (*Event, error)
}
