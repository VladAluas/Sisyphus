package extractors

import "github.com/VladAluas/Sisyphus/internal/domain"

func NewDefaultRegistry() *domain.Registry {
	r := domain.NewRegistry()

	r.Register("Postgresql", &PostgresqlExtractor{})
	r.Register("CSV", &CsvExtractor{})
	r.Register("API", &APIExtractor{})

	return r
}
