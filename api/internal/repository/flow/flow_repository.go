package flow

import (
	"context"
	"errors"
	"time"

	"ava/api/internal/logger"
	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *flowRepository) CreateFlowWithSteps(ctx context.Context, flow *model.Flow, steps []model.FlowStep) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Create(flow).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrFlowAlreadyExists
			}

			return err
		}

		if len(steps) == 0 {
			return nil
		}

		for i := range steps {
			steps[i].TenantID = flow.TenantID
			steps[i].FlowID = flow.ID
		}

		if err := dbTx.Create(&steps).Error; err != nil {
			return err
		}

		logger.Info("flowRepository.CreateFlowWithSteps",
			logger.Any("flow.ID", flow.ID),
			logger.Int("steps", len(steps)),
		)

		return nil
	})
}

func (r *flowRepository) GetByUserAndType(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*model.Flow, error) {
	var flow model.Flow

	err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Where("tenant_id = ?", tenantID).Order("step_order ASC")
		}).
		Where("tenant_id = ? AND user_id = ? AND flow_type = ?", tenantID, userID, flowType).
		First(&flow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrFlowNotFound
	}

	if err != nil {
		return nil, err
	}

	return &flow, nil
}

func (r *flowRepository) UpdateFlow(ctx context.Context, tenantID uuid.UUID, flow *model.Flow) error {
	result := r.db.WithContext(ctx).Model(&model.Flow{}).
		Where("tenant_id = ? AND id = ?", tenantID, flow.ID).
		Updates(map[string]any{
			"status":          flow.Status,
			"current_step_id": flow.CurrentStepID,
			"updated_at":      time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrFlowNotFound
	}

	return nil
}

func (r *flowRepository) UpdateStep(ctx context.Context, tenantID uuid.UUID, step *model.FlowStep) error {
	result := r.db.WithContext(ctx).Model(&model.FlowStep{}).
		Where("tenant_id = ? AND id = ?", tenantID, step.ID).
		Updates(map[string]any{
			"status":     step.Status,
			"data":       step.Data,
			"errors":     step.Errors,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrStepNotFound
	}

	return nil
}
