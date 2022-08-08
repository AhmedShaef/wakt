package task

import (
	"time"
	"unsafe"

	"github.com/AhmedShaef/wakt/business/core/task/db"
)

// Task represents an individual task.
type Task struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Pid         string        `json:"pid"`
	Wid         string        `json:"wid"`
	UID         string        `json:"uid"`
	Estimated   time.Duration `json:"estimated_seconds"`
	Active      bool          `json:"active"`
	DateCreated time.Time     `json:"date_created"`
	DateUpdated time.Time     `json:"date_updated"`
	Tracked     time.Duration `json:"tracked_seconds"`
}

// NewTask contains information needed to create a new task.
type NewTask struct {
	Name             string        `json:"name" validate:"required"`
	Pid              string        `json:"pid" validate:"required"`
	Wid              string        `json:"wid"`
	UID              string        `json:"uid"`
	EstimatedSeconds time.Duration `json:"estimated_seconds"`
	TrackedSeconds   time.Duration `json:"tracked_seconds"`
}

// UpdateTask defines what information may be provided to modify an existing
// task. All fields are optional so tasks can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateTask struct {
	Name             *string        `json:"name"`
	EstimatedSeconds *time.Duration `json:"estimated_seconds"`
	Active           *bool          `json:"active"`
	TrackedSeconds   *time.Duration `json:"tracked_seconds"`
}

// =============================================================================

func toTask(dbtask db.Task) Task {
	pu := (*Task)(unsafe.Pointer(&dbtask))
	return *pu
}

func toTasksSlice(dbtask []db.Task) []Task {
	tasks := make([]Task, len(dbtask))
	for i, dbtask := range dbtask {
		tasks[i] = toTask(dbtask)
	}
	return tasks
}
