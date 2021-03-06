package wc

import (
	"context"
	"time"
)

const (
	timeToEscalate  = 10 * time.Minute
	monitorInterval = 5 * time.Second
)

var escalationLevels = [...]string{
	0: "None",
	1: "Immediate Supervisors",
	2: "Managers",
	3: "Directors",
	4: "Executives",
}

var statusLevels = [...]string{
	0: "Green",
	1: "Yellow",
	2: "Red",
}

// Workcenter is a type that holds the information relevant to the Andon status
// of a single workcenter.
type Workcenter struct {
	ID              int
	Name            string
	CurrentProduct  string
	Status          int
	EscalationLevel int
	StatusSetAt     time.Time
}

// TimeAtStatus returns the amount of time that the workcenter has been at
// the current status.
func (wc Workcenter) TimeAtStatus() time.Duration {
	return time.Now().Sub(wc.StatusSetAt)
}

// TimeTillEscalation returns the amount of time before the next escalation level
// is set.
func (wc Workcenter) TimeTillEscalation() time.Duration {
	if wc.Status == 0 {
		return 0
	}
	return timeToEscalate - wc.TimeAtStatus()
}

// StatusDescription returns the description of the work center's status.
func (wc Workcenter) StatusDescription() string {
	return statusLevels[wc.Status]
}

// EscalationLevelDescription returns the description of the work center's esclation level.
func (wc Workcenter) EscalationLevelDescription() string {
	return escalationLevels[wc.EscalationLevel]
}

//Get WorkCenter by ID
func GetWorkcenter(ctx context.Context, id int) (Workcenter, error) {
	return findOne(ctx, id)
}

// GetAllWorkcenters retrieves all of the workcenters.
func GetAllWorkcenters(ctx context.Context) ([]Workcenter, error) {
	return findAll(ctx)
}

func escalate(ctx context.Context, id int) error {
	wc, err := findOne(ctx, id)
	if err != nil {
		return err
	}
	if wc.EscalationLevel <= len(escalationLevels)-2 {
		err = updateEscalationLevel(ctx, id, wc.EscalationLevel+1)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetWorkcenterStatus sets the workcenter with the provided ID's status.
func SetWorkcenterStatus(ctx context.Context, id, status int) error {
	escalationLevel := 0
	if status != 0 {
		escalationLevel = 1
	}
	return updateStatus(ctx, id, status, escalationLevel)
}
