// Package domain is intended to contain the custom types for the app for more flexibility
package domain

type ExecutionPlan struct {
	Batch  Batch
	Unit []ExecutionUnit
}

type ExecutionUnit struct {
  Layer   Layer
	Units []Module
}
