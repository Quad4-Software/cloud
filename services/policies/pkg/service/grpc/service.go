package grpcSVC

import (
	"context"

	v0 "github.com/Quad4-Software/cloud/protogen/gen/quad4/services/policies/v0"
	"github.com/Quad4-Software/cloud/services/policies/pkg/engine"
	"github.com/Quad4-Software/cloud/services/policies/pkg/metrics"
)

// Service defines the service handlers.
type Service struct {
	engine  engine.Engine
	metrics *metrics.Metrics
}

// New returns a service implementation for Service.
func New(engine engine.Engine, metrics *metrics.Metrics) (Service, error) {
	svc := Service{
		engine:  engine,
		metrics: metrics,
	}

	return svc, nil
}

// Evaluate exposes the engine policy evaluation.
func (s Service) Evaluate(ctx context.Context, request *v0.EvaluateRequest, response *v0.EvaluateResponse) error {
	s.metrics.EventsReceived.WithLabelValues().Inc()

	env, err := engine.NewEnvironmentFromPB(request.Environment)
	if err != nil {
		s.metrics.FailedGrpcEvaluations.WithLabelValues().Inc()
		return err
	}

	result, err := s.engine.Evaluate(ctx, request.Query, env)
	response.Result = result

	if err != nil {
		s.metrics.FailedGrpcEvaluations.WithLabelValues().Inc()
	} else {
		if result {
			s.metrics.GrpcEvaluationsThatAllow.WithLabelValues().Inc()
		} else {
			s.metrics.GrpcEvaluationsThatDontAllow.WithLabelValues().Inc()
		}
	}

	return err
}
