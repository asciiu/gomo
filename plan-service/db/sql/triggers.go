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
			order_id = id`,
		sql.Named("id", orderID))

	return err
}

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
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			uuid.FromStringOrNil(trigger.TriggerID),
			uuid.FromStringOrNil(trigger.OrderID),
			trigger.TriggerNumber,
			trigger.TriggerTemplateID,
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
	if err != nil {
		return err
	}
	return nil
}
