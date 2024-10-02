package api

import (
	"context"
	"errors"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/entity"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-devices/internal/errs"
)

const (
	dtoStatusOn  = "on"
	dtoStatusOff = "off"
)

const (
	errorCodeInternal       = "internal_code"
	errorCodeInvalidRequest = "invalid_request"
	errorCodeNotFound       = "not_found"
)

type service interface {
	RegisterDevice(ctx context.Context, deviceID int64) error
	SetStatus(ctx context.Context, in entity.SetStatus) error
	GetStatus(ctx context.Context, deviceID int64) (entity.DeviceSettings, error)
}

type server struct {
	service service
}

func NewServer(service service) *server {
	return &server{service: service}
}

func (s *server) GetDevicesDeviceIdStatus(ctx context.Context, request GetDevicesDeviceIdStatusRequestObject) (GetDevicesDeviceIdStatusResponseObject, error) {
	if !isValidDeviceID(request.DeviceId) {
		return GetDevicesDeviceIdStatus400JSONResponse{
			Code:    errorCodeInvalidRequest,
			Message: "invalid device_id",
		}, nil
	}

	settings, err := s.service.GetStatus(ctx, int64(request.DeviceId))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return GetDevicesDeviceIdStatus404JSONResponse{
				Code:    errorCodeNotFound,
				Message: err.Error(),
			}, nil
		}
		return GetDevicesDeviceIdStatus500JSONResponse{
			Code:    errorCodeInternal,
			Message: err.Error(),
		}, nil
	}

	return convertSettingsToOutDto(settings), nil
}

func (s *server) PutDevicesDeviceIdStatus(ctx context.Context, request PutDevicesDeviceIdStatusRequestObject) (PutDevicesDeviceIdStatusResponseObject, error) {
	if request.Body == nil {
		return PutDevicesDeviceIdStatus400JSONResponse{
			Code:    errorCodeInvalidRequest,
			Message: "nil body",
		}, nil
	}

	if !isValidStatus(request.Body.Status) {
		return PutDevicesDeviceIdStatus400JSONResponse{
			Code:    errorCodeInvalidRequest,
			Message: "invalid status",
		}, nil
	}

	if !isValidDeviceID(request.DeviceId) {
		return PutDevicesDeviceIdStatus400JSONResponse{
			Code:    errorCodeInvalidRequest,
			Message: "invalid device_id",
		}, nil
	}

	err := s.service.SetStatus(ctx, convertInSettingsToEntity(request.DeviceId, *request.Body))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return PutDevicesDeviceIdStatus404JSONResponse{
				Code:    errorCodeNotFound,
				Message: err.Error(),
			}, nil
		}
		return PutDevicesDeviceIdStatus500JSONResponse{
			Code:    errorCodeInternal,
			Message: err.Error(),
		}, nil
	}
	return PutDevicesDeviceIdStatus204Response{}, nil
}

func (s *server) PostDevices(ctx context.Context, request PostDevicesRequestObject) (PostDevicesResponseObject, error) {
	if request.Body == nil {
		return PostDevices400JSONResponse{
			Code:    errorCodeInvalidRequest,
			Message: "nil body",
		}, nil
	}

	if !isValidDeviceID(request.Body.DeviceId) {
		return PostDevices400JSONResponse{
			Code:    errorCodeInvalidRequest,
			Message: "invalid device_id",
		}, nil
	}

	err := s.service.RegisterDevice(ctx, int64(request.Body.DeviceId))
	if err != nil {
		return PostDevices500JSONResponse{
			Code:    errorCodeInternal,
			Message: err.Error(),
		}, nil
	}
	return PostDevices201Response{}, nil
}

func convertSettingsToOutDto(settings entity.DeviceSettings) GetDevicesDeviceIdStatus200JSONResponse {
	status := dtoStatusOn
	if !settings.Enabled {
		status = dtoStatusOff
	}
	return GetDevicesDeviceIdStatus200JSONResponse{
		DeviceId:  int(settings.DeviceID),
		Status:    status,
		UpdatedAt: &settings.UpdatedAt,
	}
}

func convertInSettingsToEntity(deviceID int, in PutDevicesDeviceIdStatusJSONRequestBody) entity.SetStatus {
	e := entity.SetStatus{
		DeviceID: int64(deviceID),
	}
	if in.Status == dtoStatusOn {
		e.Enabled = true
	}
	return e
}

func isValidStatus(status string) bool {
	return status == dtoStatusOn || status == dtoStatusOff
}

func isValidDeviceID(deviceID int) bool {
	return deviceID > 0
}
