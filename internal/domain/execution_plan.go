// Package domain is intended to contain the custom types for the app for more flexibility
package domain

type ExecutionPlan struct {
	Batch  Batch
	Unit []ExecutionLayer
}

type ExecutionLayer struct {
  Layer   Layer
	Units []ExecutionUnit
}

type ExecutionUnit struct {
	Module Module
}
