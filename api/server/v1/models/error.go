package models

type ResponseType string

const (
	// ResponseTypeSuccess captures enum value "success"
	ResponseTypeSuccess ResponseType = "success"
	// ResponseTypeBusinessConflict captures enum value "business_conflict"
	ResponseTypeBusinessConflict ResponseType = "business_conflict"
	// ResponseTypeUnprocessableEntity captures enum value "unprocessable_entity"
	ResponseTypeUnprocessableEntity ResponseType = "unprocessable_entity"
	// ResponseTypeBadParameters captures enum value "bad_parameters"
	ResponseTypeBadParameters ResponseType = "bad_parameters"
	// ResponseTypeInternalError captures enum value "internal_error"
	ResponseTypeInternalError ResponseType = "internal_error"
	// ResponseTypeNotFound captures enum value "not_found"
	ResponseTypeNotFound ResponseType = "not_found"
	// ResponseTypeSecurityError captures enum value "security_error"
	ResponseTypeSecurityError ResponseType = "security_error"
	// ResponseTypePermissionError captures enum value "permission_error"
	ResponseTypePermissionError ResponseType = "permission_error"
)

type ErrorErrorsItems struct {

	// тип ишибки
	Code string `json:"code,omitempty"`

	// поле вызвавшее ошибку
	Field string `json:"field,omitempty"`

	// описание
	Message string `json:"message,omitempty"`
}


type ErrorErrors []*ErrorErrorsItems

// swagger:model
type Error struct {

	// code
	Code ResponseType `json:"code,omitempty"`

	// errors
	Errors ErrorErrors `json:"errors"`

	// описание ошибки
	Message string `json:"message,omitempty"`
}

