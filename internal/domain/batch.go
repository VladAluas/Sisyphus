// Package domain is intended to contain the custom types for the app for more flexibility
package domain

type Batch struct {
	BatchID         string
	BatchCode       string
	BatchRunID      string
	BatchParameters map[string]string
}
