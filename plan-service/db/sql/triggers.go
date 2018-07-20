package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/lib/pq"
)

func InsertTriggers(db *sql.DB, triggers []*protoOrder.Trigger) error {
	totalColumns := 9
	valueStrings := make([]string, 0, len(triggers))
	valueArgs := make([]interface{}, 0, len(triggers)*totalColumns)

	for _, trigger := range triggers {
		jsonCode, err := json.Marshal(trigger.Code)
		if err != nil {
			return err
		}

		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, trigger.TriggerID)
		valueArgs = append(valueArgs, trigger.OrderID)
		valueArgs = append(valueArgs, trigger.TriggerNumber)
		valueArgs = append(valueArgs, trigger.TriggerTemplateID)
		valueArgs = append(valueArgs, trigger.Name)
		valueArgs = append(valueArgs, jsonCode)
		valueArgs = append(valueArgs, pq.Array(trigger.Actions))
		valueArgs = append(valueArgs, trigger.CreatedOn)
		valueArgs = append(valueArgs, trigger.UpdatedOn)
	}
	stmt := fmt.Sprintf(`INSERT INTO triggers (
		id, 
		order_id, 
		trigger_number,
		trigger_template_id,
		name,
		code,
		actions,
		created_on,
		updated_on) VALUES %s`, strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}
