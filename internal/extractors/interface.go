package extractors

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type Extractor interface {
	Extract(ctx context.Context, unit domain.ExecutionUnit) error
}
