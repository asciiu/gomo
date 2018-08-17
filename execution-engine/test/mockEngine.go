package test

import (
	"context"
	"database/sql"

	protoEngine "github.com/asciiu/gomo/execution-engine/proto/engine"
	"github.com/micro/go-micro/client"
)

// Test clients of the Key service should use this client interface.
type mockEngine struct {
	db    *sql.DB
	Plans []*protoEngine.Plan
}

func (m *mockEngine) AddPlan(ctx context.Context, in *protoEngine.NewPlanRequest, opts ...client.CallOption) (*protoEngine.PlanResponse, error) {
	plan := protoEngine.Plan{
		PlanID: in.PlanID,
	}
	m.Plans = append(m.Plans, &plan)

	return &protoEngine.PlanResponse{
		Status: "success",
		Data: &protoEngine.PlanList{
			Plans: []*protoEngine.Plan{&plan},
		},
	}, nil
}

func (m *mockEngine) GetActivePlans(ctx context.Context, in *protoEngine.ActiveRequest, opts ...client.CallOption) (*protoEngine.PlanResponse, error) {
	plans := make([]*protoEngine.Plan, 0)
	return &protoEngine.PlanResponse{
		Status: "success",
		Data: &protoEngine.PlanList{
			Plans: plans,
		},
	}, nil
}

func (m *mockEngine) KillPlan(ctx context.Context, in *protoEngine.KillRequest, opts ...client.CallOption) (*protoEngine.PlanResponse, error) {
	plans := make([]*protoEngine.Plan, 0)
	for _, p := range m.Plans {
		if p.PlanID == in.PlanID {
			plans = append(plans,
				&protoEngine.Plan{
					PlanID: in.PlanID,
				},
			)
		}
	}

	return &protoEngine.PlanResponse{
		Status: "success",
		Data: &protoEngine.PlanList{
			Plans: plans,
		},
	}, nil
}

func (m *mockEngine) KillUserPlans(ctx context.Context, in *protoEngine.KillUserRequest, opts ...client.CallOption) (*protoEngine.PlanResponse, error) {
	plans := make([]*protoEngine.Plan, 0)
	return &protoEngine.PlanResponse{
		Status: "success",
		Data: &protoEngine.PlanList{
			Plans: plans,
		},
	}, nil
}

func MockEngineClient(db *sql.DB) protoEngine.ExecutionEngineClient {
	return &mockEngine{
		db:    db,
		Plans: make([]*protoEngine.Plan, 0),
	}
}
