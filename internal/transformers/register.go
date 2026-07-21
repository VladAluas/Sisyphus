// Package transformers focuses mainly on data cleaning
package transformers

func NewCleanRegistry() *CleanersRegistry {
	r := NewRegistry()

	r.Register("Clean", &GeneralCleaner{})

	return r
}
