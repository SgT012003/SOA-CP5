package errors

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Name    string `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, name, message string) *AppError {
	return &AppError{
		Code:    code,
		Name:    name,
		Message: message,
	}
}

var (
	ErrInvalidDateRange    = NewAppError(http.StatusBadRequest, "InvalidDateRangeException", "A data de check-out deve ser maior que a data de check-in.")
	ErrRoomUnavailable     = NewAppError(http.StatusConflict, "RoomUnavailableException", "Quarto não disponível no período selecionado.")
	ErrCapacityExceeded    = NewAppError(http.StatusBadRequest, "CapacityExceededException", "A capacidade do quarto foi excedida.")
	ErrInvalidState        = NewAppError(http.StatusConflict, "InvalidReservationStateException", "Transição de estado da reserva inválida.")
	ErrRoomCannotBeDeleted = NewAppError(http.StatusConflict, "RoomCannotBeDeletedException", "Quartos não podem ser deletados, apenas inativados.")
	ErrNotFound            = NewAppError(http.StatusNotFound, "NotFoundException", "Recurso não encontrado.")
	ErrValidation          = NewAppError(http.StatusBadRequest, "ValidationException", "Erro de validação nos dados enviados.")
	ErrInternalServer      = NewAppError(http.StatusInternalServerError, "InternalServerException", "Erro interno no servidor.")
)
