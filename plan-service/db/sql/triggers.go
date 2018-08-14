package sql

import (
	"context"
	"database/sql"

	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

func DeleteTriggersWithOrderID(txn *sql.Tx, ctx context.Context, orderID string) error {
	_, err := txn.ExecContext(ctx, `
		DELETE FROM triggers 
		WHERE
			order_id = $1`, orderID)

	return err
}

func InsertTriggers(txn *sql.Tx, triggers []*protoOrder.Trigger) error {
	stmt, err := txn.Prepare(pq.CopyIn("triggers",
		"id",
		"order_id",
		"trigger_template_id",
		"index",
		"title",
		"name",
		"code",
		"actions",
		"created_on",
		"updated_on"))

	if err != nil {
		return err
	}

	for _, trigger := range triggers {
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			uuid.FromStringOrNil(trigger.TriggerID),
			uuid.FromStringOrNil(trigger.OrderID),
			trigger.TriggerTemplateID,
			trigger.Index,
			trigger.Title,
			trigger.Name,
			trigger.Code,
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
	return err
}

func UpdateTriggerResults(db *sql.DB, triggerID string, triggeredPrice float64, triggeredCondition, triggeredTimestamp string) error {
	updatePlanOrderSql := `
	UPDATE triggers 
	SET 
		triggered_price = $1,
		triggered_condition = $2,
		triggered_timestamp = $3,
	WHERE id = $4`

	_, err := db.Exec(updatePlanOrderSql,
		triggeredPrice,
		triggeredCondition,
		triggeredTimestamp,
		triggerID)

	return err
}
