package sql

import (
	"database/sql"
	"encoding/json"

	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

// func InsertTriggers(db *sql.DB, triggers []*protoOrder.Trigger) error {
// 	totalColumns := 9
// 	valueStrings := make([]string, 0, len(triggers))
// 	valueArgs := make([]interface{}, 0, len(triggers)*totalColumns)

// 	for _, trigger := range triggers {
// 		jsonCode, err := json.Marshal(trigger.Code)
// 		if err != nil {
// 			return err
// 		}

// 		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
// 		valueArgs = append(valueArgs, trigger.TriggerID)
// 		valueArgs = append(valueArgs, trigger.OrderID)
// 		valueArgs = append(valueArgs, trigger.TriggerNumber)
// 		valueArgs = append(valueArgs, trigger.TriggerTemplateID)
// 		valueArgs = append(valueArgs, trigger.Name)
// 		valueArgs = append(valueArgs, jsonCode)
// 		valueArgs = append(valueArgs, pq.Array(trigger.Actions))
// 		valueArgs = append(valueArgs, trigger.CreatedOn)
// 		valueArgs = append(valueArgs, trigger.UpdatedOn)
// 	}
// 	stmt := fmt.Sprintf(`INSERT INTO triggers (
// 		id,
// 		order_id,
// 		trigger_number,
// 		trigger_template_id,
// 		name,
// 		code,
// 		actions,
// 		created_on,
// 		updated_on) VALUES %s`, strings.Join(valueStrings, ","))
// 	_, err := db.Exec(stmt, valueArgs...)

// 	if err != nil {
// 		return errors.New("bulk insert InsertTriggers failed: " + err.Error())
// 	}

// 	return nil
// }

func InsertTriggers(txn *sql.Tx, triggers []*protoOrder.Trigger) error {
	stmt, err := txn.Prepare(pq.CopyIn("triggers",
		"id",
		"order_id",
		"trigger_number",
		"trigger_template_id",
		"name",
		"code",
		"actions",
		"created_on",
		"updated_on"))

	if err != nil {
		return err
	}

	for _, trigger := range triggers {
		jsonCode, err := json.Marshal(trigger.Code)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			uuid.FromStringOrNil(trigger.TriggerID),
			uuid.FromStringOrNil(trigger.OrderID),
			trigger.TriggerNumber,
			trigger.TriggerTemplateID,
			trigger.Name,
			jsonCode,
			pq.Array(trigger.Actions),
			trigger.CreatedOn,
			trigger.UpdatedOn)

		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}
