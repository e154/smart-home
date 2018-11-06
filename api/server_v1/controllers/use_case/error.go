package use_case

import (
	stubMoldes "github.com/e154/smart-home/api/server_v1/stub/models"
	"fmt"
	"github.com/go-openapi/runtime"
	"net/http"
	"github.com/e154/smart-home/system/validation"
	"github.com/iancoleman/strcase"
)

const (
	FieldNotValid      = "field_not_valid"
	FieldNotBlank      = "field_not_blank"
	FieldSizeMax       = "field_size_max"
	FieldSizeMin       = "field_size_min"
	FieldInvalidLength = "field_invalid_length"
	FieldNotValidChars = "field_not_valid_chars"
	FieldMax           = "field_max"
	FieldMin           = "field_min"
	FieldFuture        = "field_future"
	FieldPast          = "field_past"
	FieldEmail         = "field_email"
	FieldCardNumber    = "field_card_number"
	FieldPhone         = "field_phone"
	FieldDuplicate     = "field_duplicate"
)

const (
	UserBlocked         = "user_blocked"
	NewCodeGenerated    = "new_code_generated"
	CodeWrong           = "code_wrong"
	DetectSpam          = "detect spam"
	RefreshTokenInvalid = "refresh token invalid"
	AccessTokenInvalid  = "access token invalid"
	NeedUpdateWarning   = "need update warning"
	NeedUpdateCritical  = "need update critical"
)

type Error struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *stubMoldes.ErrorModel `json:"body,omitempty"`
}

// NewPostAuthTokenDefault creates Error with default headers values
func NewError(code int, msg ...interface{}) *Error {

	var _code, _message string

	switch code {
	case 400:
		_code = "bad_parameters"
	case 401:
		_code = "security_error"
	case 403:
		_code = "permission_error"
	case 404:
		_code = "not_found"
	case 409:
		_code = "business_conflict"
	case 422:
		_code = "unprocessable_entity"
	case 426:
		_code = "need_update"
	case 500:
		_code = "internal_error"
	default:
		_code = "internal_error"
	}

	_message = fmt.Sprintf("general code for %d http code", code)

	if len(msg) > 0 {
		switch v := msg[0].(type) {
		case error:
			_message = v.Error()
		case string:
			_message = v
		}
	}

	return &Error{
		_statusCode: code,
		Payload: &stubMoldes.ErrorModel{
			Code:    stubMoldes.ResponseType(_code),
			Message: _message,
			Errors:  stubMoldes.ErrorModelErrors{},
		},
	}
}

// WriteResponse to the client
func (o *Error) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	//log.Errorf("%s: %s", o.Payload.Code, o.Payload.Message)
	if len(o.Payload.Errors) > 0 {
		for _, err := range o.Payload.Errors {
			log.Errorf("%s: %s", err.Field, err.Message)
		}
	}

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			//panic(err)
		}
	}
}

func (o *Error) Fields() []*stubMoldes.ErrorModelErrorsItems {
	return o.Payload.Errors
}

func (o *Error) AddField(code, message, field string) *Error {

	_field := &stubMoldes.ErrorModelErrorsItems{
		Code:    code,
		Message: message,
		Field:   field,
	}

	o.Payload.Errors = append(o.Payload.Errors, _field)

	return o
}

func (o *Error) SetMessage(err error) *Error {
	o.Payload.Message = err.Error()
	return o
}

func (o *Error) CheckNum(num interface{}, name string, min, max float64) *Error {

	var n float64

	switch i := num.(type) {
	case float64:
		n = i
	case int64:
		n = float64(i)
	case int32:
		n = float64(i)
	case int:
		n = float64(i)
	case float32:
		n = float64(i)
	default:
		return o.AddFieldf("", FieldNotValid)
	}

	if n < min {
		o.AddField(fmt.Sprintf("common.%s_min", name), fmt.Sprintf("The %s can't be less than %v", name, min), name)
	}
	if n > max {
		o.AddField(fmt.Sprintf("common.%s_max", name), fmt.Sprintf("The %s can't be greater than %v", name, max), name)
	}

	return o
}

func (o *Error) CheckString(s, name string, min, max int) *Error {

	return o
}

func (o *Error) AddFieldf(name, code string, N ...int) *Error {

	var n int
	if len(N) > 0 {
		n = N[0]
	}

	switch code {
	case FieldNotValid:
		o.AddField("common.field_not_valid", fmt.Sprintf("The %s isn't valid", name), name)
	case FieldNotBlank:
		o.AddField("common.field_not_blank", "The field can't be empty", name)
	case FieldSizeMax:
		o.AddField("common.field_size_max", "The field is too long", name)
	case FieldSizeMin:
		o.AddField("common.field_size_min", "The field is too short", name)
	case FieldInvalidLength:
		o.AddField("common.field_invalid_length", "The field length is not correct", name)
	case FieldNotValidChars:
		o.AddField("common.field_not_valid_chars", "The field contains invalid characters", name)
	case FieldMax:
		o.AddField("common.field_max", fmt.Sprintf("The nuber can't be greater than %d", n), name)
	case FieldMin:
		o.AddField("common.field_min", fmt.Sprintf("The nuber can't be less than %d", n), name)
	case FieldFuture:
		o.AddField("common.field_future", fmt.Sprintf("The date should be later than %d", n), name)
	case FieldPast:
		o.AddField("common.field_past", fmt.Sprintf("The date should be early than %d", n), name)
	case FieldEmail:
		o.AddField("common.field_email", "Email isn't valid", name)
	case FieldCardNumber:
		o.AddField("common.field_card_number", "Card number isn't valid", name)
	case FieldPhone:
		o.AddField("common.field_phone", "The phone number isn't valid", name)
	case FieldDuplicate:
		o.AddField("common.field_duplicate", "The field value should be unique", name)
	case UserBlocked:
		o.AddField("sec.user_blocked", "User is blocked", name)
	case NewCodeGenerated:
		o.AddField("sec.new_code_generated", "New SMS code generated", name)
	case CodeWrong:
		o.AddField("sec.code_wrong", "Code is wrong", name)
	case DetectSpam:
		o.AddField("sec.detect_spam", "Detect spam", name)
	case RefreshTokenInvalid:
		o.AddField("sec.refresh_token_invalid", "Refresh token is invalid", name)
	case AccessTokenInvalid:
		o.AddField("sec.access_token_invalid", "Access token is invalid", name)
	case NeedUpdateWarning:
		o.AddField("sec.need_update_warning", "Need update app warning", name)
	case NeedUpdateCritical:
		o.AddField("sec.need_update_critical", "Need update app critical", name)
	}

	return o
}

func (o *Error) Errors() bool {
	return len(o.Payload.Errors) > 0
}

func (e *Error) Error() string {
	fmt.Println(e.Payload.Errors[0].Message)
	return e.Payload.Message
}

func (e *Error) ValidationToErrors(errs []*validation.Error) *Error {


	for _, err := range errs {
		field := strcase.ToSnake(err.Field)

		var code string
		var limitValue []int

		switch err.Name {
		case "Required":
			code = FieldNotBlank
		case "Match":
			code = FieldNotValid
		case "MinSize":
			code = FieldSizeMin
		case "MaxSize":
			code = FieldSizeMax
		case "Email":
			code = FieldEmail
		case "Min":
			code = FieldMin
		case "Max":
			code = FieldMax
		default:
			log.Warningf("не известный тип валидации: %s", err.Name)
		}

		limit, ok := err.LimitValue.(int)
		if ok {
			limitValue = append(limitValue, limit)
		}

		if code != "" {
			e.AddFieldf(field, code, limitValue...)
		}
	}

	return e
}