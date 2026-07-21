// Package transformers focuses mainly on data cleaning
package transformers

import "github.com/VladAluas/Sisyphus/internal/domain"

func NewDefaultRegistry() *domain.Registry {
	r := domain.NewRegistry()

	r.Register("Clean", &GeneralCleaner{})

	return r
}

