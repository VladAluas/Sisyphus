package extractors

func NewExtractorsRegistry() *ExtractorsRegistry {
	r := NewRegistry()

	r.Register("Postgresql", &PostgresqlExtractor{})
	r.Register("CSV", &CsvExtractor{})
	r.Register("API", &APIExtractor{})

	return r
}
