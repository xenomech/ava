package flow

import (
	"context"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateFlowWithSteps(ctx context.Context, flow *model.Flow, steps []model.FlowStep) error
	GetByUserAndType(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*model.Flow, error)
	UpdateFlow(ctx context.Context, tenantID uuid.UUID, flow *model.Flow) error
	UpdateStep(ctx context.Context, tenantID uuid.UUID, step *model.FlowStep) error
	ReplaceSteps(ctx context.Context, tenantID uuid.UUID, flow *model.Flow, steps []model.FlowStep) error
}

type flowRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &flowRepository{db: db}
}
