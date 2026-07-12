// Package domain is intended to contain the custom types for the app for more flexibility
package domain

type ExecutionPlan struct {
	Batch  Batch
	Layers []ExecutionLayer
}

type ExecutionLayer struct {
  Layer   Layer
	Status  string
	Modules []ExecutionUnit
}

type ExecutionUnit struct {
	Module Module
	Source Source
	Status string
}
