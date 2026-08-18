package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StepErrors map[string]string

func (errs StepErrors) Value() (driver.Value, error) {
	if len(errs) == 0 {
		return []byte("{}"), nil
	}

	return json.Marshal(errs)
}

func (errs *StepErrors) Scan(value any) error {
	if value == nil {
		*errs = StepErrors{}

		return nil
	}

	var raw []byte

	switch typed := value.(type) {
	case []byte:
		raw = typed
	case string:
		raw = []byte(typed)
	default:
		return fmt.Errorf("unsupported type for StepErrors: %T", value)
	}

	if len(raw) == 0 {
		*errs = StepErrors{}

		return nil
	}

	return json.Unmarshal(raw, errs)
}

type FlowStatus string

const (
	FlowStatusPending    FlowStatus = "pending"
	FlowStatusInProgress FlowStatus = "in_progress"
	FlowStatusCompleted  FlowStatus = "completed"
	FlowStatusFailed     FlowStatus = "failed"
)

type FlowStepStatus string

const (
	FlowStepStatusPending    FlowStepStatus = "pending"
	FlowStepStatusInProgress FlowStepStatus = "in_progress"
	FlowStepStatusCompleted  FlowStepStatus = "completed"
	FlowStepStatusSkipped    FlowStepStatus = "skipped"
	FlowStepStatusFailed     FlowStepStatus = "failed"
)

type Flow struct {
	BaseModel
	TenantID      uuid.UUID       `gorm:"type:uuid;not null;index:idx_flow_tenant_user_type,unique" json:"tenant_id"`
	Tenant        *Tenant         `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	UserID        uuid.UUID       `gorm:"type:uuid;not null;index:idx_flow_tenant_user_type,unique" json:"user_id"`
	User          *User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	FlowType      string          `gorm:"type:varchar(50);not null;index:idx_flow_tenant_user_type,unique" json:"flow_type"`
	Status        FlowStatus      `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CurrentStepID string          `gorm:"type:varchar(50);not null" json:"current_step_id"`
	Metadata      json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	Steps         []FlowStep      `gorm:"foreignKey:FlowID" json:"steps"`
}

func (flow *Flow) BeforeCreate(tx *gorm.DB) error {
	return flow.BaseModel.BeforeCreate(tx)
}

func NewFlow(tenantID, userID uuid.UUID, flowType, currentStepID string, metadata json.RawMessage) *Flow {
	return &Flow{
		TenantID:      tenantID,
		UserID:        userID,
		FlowType:      flowType,
		Status:        FlowStatusInProgress,
		CurrentStepID: currentStepID,
		Metadata:      metadata,
	}
}

type FlowStep struct {
	BaseModel
	TenantID    uuid.UUID       `gorm:"type:uuid;not null;index:idx_flow_step_tenant_flow" json:"tenant_id"`
	Tenant      *Tenant         `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	FlowID      uuid.UUID       `gorm:"type:uuid;not null;index;index:idx_flow_step_tenant_flow" json:"flow_id"`
	StepID      string          `gorm:"type:varchar(50);not null" json:"step_id"`
	Title       string          `gorm:"type:varchar(100);not null" json:"title"`
	Description string          `gorm:"type:text" json:"description"`
	StepOrder   int             `gorm:"not null" json:"order"`
	Status      FlowStepStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Skippable   bool            `gorm:"default:false" json:"skippable"`
	Data        json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"data"`
	Errors      StepErrors      `gorm:"type:jsonb;default:'{}'" json:"errors"`
}

func (step *FlowStep) BeforeCreate(tx *gorm.DB) error {
	return step.BaseModel.BeforeCreate(tx)
}

func NewFlowStep(tenantID, flowID uuid.UUID, stepID, title, description string, order int, status FlowStepStatus, skippable bool) FlowStep {
	return FlowStep{
		TenantID:    tenantID,
		FlowID:      flowID,
		StepID:      stepID,
		Title:       title,
		Description: description,
		StepOrder:   order,
		Status:      status,
		Skippable:   skippable,
		Data:        json.RawMessage("{}"),
		Errors:      StepErrors{},
	}
}
