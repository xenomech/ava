package flow

import "ava/api/internal/serrors"

var (
	ErrInvalidFlowType      = serrors.NewCoded("invalid_flow_type", "invalid flow type")
	ErrFlowAlreadyCompleted = serrors.NewCoded("flow_already_completed", "flow already completed")
	ErrStepNotFound         = serrors.NewCoded("step_not_found", "step not found")
	ErrStepNotCurrent       = serrors.NewCoded("step_not_current", "step is not the current step")
	ErrStepNotSkippable     = serrors.NewCoded("step_not_skippable", "step is not skippable")
	ErrNoPreviousStep       = serrors.NewCoded("no_previous_step", "no previous step")
	ErrStepValidationFailed = serrors.NewCoded("step_validation_failed", "step validation failed")
	ErrInvalidStepData      = serrors.NewCoded("invalid_step_data", "step data does not match the shape this step expects")
	ErrStepNotPermitted     = serrors.NewCoded("step_not_permitted", "not permitted to complete this step")
)
