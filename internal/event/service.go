package event

import (
	"context"

	"github.com/krissukoco/go-event-ticket/internal/models"
)

type Service interface {
	GetById(ctx context.Context, id string) (*models.Event, error)
	GetAll(ctx context.Context, page, limit int) ([]*models.Event, error)
	Insert(ctx context.Context, in *models.Event) (*models.Event, error)
}
