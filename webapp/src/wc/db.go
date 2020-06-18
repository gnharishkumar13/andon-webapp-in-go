package wc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/user/andon-webapp-in-go/src/db"
)

var database *sql.DB

func SetDB(db *sql.DB) {
	database = db
}

func findOne(ctx context.Context, id int) (Workcenter, error) {
	result := database.QueryRowContext(ctx,
		`SELECT 
			id, wc_name, current_product, wc_status, escalation_level, status_set_at 
		FROM workcenters
		WHERE id = $1`, id)
	wc := Workcenter{}
	var statusSetAtRaw string
	err := result.Scan(&wc.ID, &wc.Name, &wc.CurrentProduct,
		&wc.Status, &wc.EscalationLevel, &statusSetAtRaw)
	if err != nil {
		return Workcenter{},
			fmt.Errorf("failed to retrieve workcenter from database with id %q: %v", id, err)
	}
	statusSetAt, err := db.FromTimestamp(statusSetAtRaw)
	if err != nil {
		return Workcenter{},
			fmt.Errorf("failed to parse status_set_at timestamp for workcenter with id %q: %v", id, err)
	}
	wc.StatusSetAt = statusSetAt
	return wc, nil
}

func findAll(ctx context.Context) ([]Workcenter, error) {
	result, err := database.QueryContext(ctx,
		`SELECT 
			id, wc_name, current_product, wc_status, escalation_level, status_set_at 
		FROM workcenters
		ORDER BY wc_name`)
	if err != nil {
		msg := fmt.Sprintf("failed to retrieve workcententers from database: %v", err)
		log.Println(msg)
		return nil, fmt.Errorf(msg)
	}
	workcenters := []Workcenter{}
	var wc Workcenter
	var statusSetAtRaw string
	for result.Next() {
		err := result.Scan(&wc.ID, &wc.Name, &wc.CurrentProduct, &wc.Status, &wc.EscalationLevel, &statusSetAtRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve workcenter fields database: %v", err)
		}
		statusSetAt, err := time.Parse(db.TimeLayoutTimestamp, statusSetAtRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse status_set_at timestamp for workcenter with id %q: %v", wc.ID, err)
		}
		wc.StatusSetAt = statusSetAt
		workcenters = append(workcenters, wc)
	}
	return workcenters, nil
}

func updateEscalationLevel(ctx context.Context, id, escalationLevel int) error {
	_, err := database.ExecContext(ctx,
		`UPDATE workcenters
		SET status_set_at=$1,
			escalation_level=$2
		WHERE id=$3
		`, db.ToTimestamp(time.Now()), escalationLevel, id)
	if err != nil {
		return fmt.Errorf("failed to update escalation level of workcenter with id %q: %v", id, err)
	}
	return nil
}

func updateStatus(ctx context.Context, id int, status int, escalationLevel int) error {
	_, err := database.ExecContext(ctx,
		`UPDATE workcenters
		SET wc_status=$1,
			status_set_at=$2,
			escalation_level=$3
		WHERE id=$4
		`, status, db.ToTimestamp(time.Now()), escalationLevel, id)
	if err != nil {
		return fmt.Errorf("failed to update status of workcenter with id %q: %v", id, err)
	}
	return nil
}
