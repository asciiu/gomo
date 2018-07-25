package sql

import (
	"database/sql"
	"encoding/json"

	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

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