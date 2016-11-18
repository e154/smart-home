package models

type RedactorGrapSettings struct {
	Position	struct{
				Top int64	`json:"top"`
				Left int64	`json:"left"`
			}			`json:"position"`
}

type RedactorConnector struct {
	Id		int64			`json:"id"`
	Start		struct{
				Object int64	`json:"object"`
				Point int64	`json:"point"`
				}		`json:"start"`
	End		struct{
				Object int64	`json:"object"`
				Point int64	`json:"point"`
				}		`json:"end"`
	flow_type	string			`json:"flow_type"`
	title		string			`json:"title"`

}

type RedactorObject struct {
	Id		int64			`json:"id"`
	Type		struct {
				Name	string	`json:"name"`
				Start	interface{}	`json:"start"`
				End	interface{}	`json:"end"`
				Status	string	`json:"status"`
				Action	string	`json:"action"`
				}		`json:"type"`
	Position	struct{
				Top int64	`json:"top"`
				Left int64	`json:"left"`
			}			`json:"position"`
	Status		string                  `json:"status"`
	Error		string			`json:"error"`
	Title		string			`json:"title"`
	Description	string                  `json:"description"`
}

type RedactorFlow struct {
	Name		string			`json:"name"`
	Description	string			`json:"description"`
	Objects		[]*RedactorObject	`json:"objects"`
	Connectors	[]*RedactorConnector	`json:"connectors"`
}