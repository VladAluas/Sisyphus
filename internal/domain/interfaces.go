package domain

import (
	"context"
)

type Processor interface {
	Process(ctx context.Context, task ExecutionUnit) error
}
