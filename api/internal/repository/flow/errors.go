package flow

import "ava/api/pkg/serrors"

var (
	ErrFlowNotFound      = serrors.New("flow not found")
	ErrFlowAlreadyExists = serrors.New("flow already exists")
	ErrStepNotFound      = serrors.New("flow step not found")
)
