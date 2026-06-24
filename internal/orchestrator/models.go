// Package orchestrator manages the tasks for the ETL
package orchestrator

type ModuleTask struct {
	BatchRunID string
	BatchID    string
	ModuleID   string
	LayerID    string
	Priority   int
}

