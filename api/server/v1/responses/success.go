package responses

// Success response
// swagger:response Success
type Success struct {
	// in:body
	Body struct {
	}
}

// Success with id response
// swagger:response NewObjectSuccess
type NewObjectSuccess struct {
	// in:body
	Body struct {
		Id int64 `json:"id"`
	}
}
