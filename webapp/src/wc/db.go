package wc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/user/andon-webapp-in-go/src/db"
)

var database *sql.DB

func SetDB(db *sql.DB){
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
