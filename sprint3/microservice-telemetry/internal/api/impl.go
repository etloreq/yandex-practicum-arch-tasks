package api

import (
	"context"
	"errors"
	"time"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/entity"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/errs"
)

const (
	errorCodeInternal       = "internal_code"
	errorCodeInvalidRequest = "invalid_request"
	errorCodeNotFound       = "not_found"
)

type service interface {
	SaveTelemetry(ctx context.Context, telemetry entity.Telemetry) error
	GetLatest(ctx context.Context, deviceID int64) (entity.Telemetry, error)
}

type server struct {
	service service
}

func NewServer(service service) *server {
	return &server{service: service}
}

func (s *server) PostTelemetry(ctx context.Context, request PostTelemetryRequestObject) (PostTelemetryResponseObject, error) {
	if errValidation := validateNewTelemetry(request); errValidation != nil {
		return PostTelemetry400JSONResponse{
			Code:    errorCodeInvalidRequest,
			Message: errValidation.Error(),
		}, nil
	}

	err := s.service.SaveTelemetry(ctx, entity.Telemetry{
		DeviceID:  int64(request.Body.DeviceId),
		Measure:   int64(request.Body.Measure),
		Timestamp: time.Unix(int64(request.Body.Timestamp), 0),
	})
	if err != nil {
		return PostTelemetry500JSONResponse{
			Code:    errorCodeInternal,
			Message: err.Error(),
		}, nil
	}

	return PostTelemetry201Response{}, nil
}

func (s *server) GetTelemetryDeviceIdLatest(ctx context.Context, request GetTelemetryDeviceIdLatestRequestObject) (GetTelemetryDeviceIdLatestResponseObject, error) {
	telemetry, err := s.service.GetLatest(ctx, int64(request.DeviceId))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return GetTelemetryDeviceIdLatest404JSONResponse{
				Code:    errorCodeNotFound,
				Message: err.Error(),
			}, nil
		}
		return GetTelemetryDeviceIdLatest500JSONResponse{
			Code:    errorCodeInternal,
			Message: err.Error(),
		}, nil
	}
	return GetTelemetryDeviceIdLatest200JSONResponse{
		DeviceId:  int(telemetry.DeviceID),
		Measure:   int(telemetry.Measure),
		Timestamp: int(telemetry.Timestamp.Unix()),
	}, nil
}

func validateNewTelemetry(request PostTelemetryRequestObject) error {
	if request.Body == nil {
		return errors.New("nil body")
	}
	if request.Body.DeviceId <= 0 {
		return errors.New("invalid device id")
	}
	return nil
}
