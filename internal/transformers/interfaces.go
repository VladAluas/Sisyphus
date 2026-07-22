package transformers

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type Cleaner interface {
	Clean(ctx context.Context, unit domain.ExecutionUnit) error
}
