// Smart home
//
// This documentation describes APIs found under https://github.com/e154/smart-home
//
//     BasePath: /api/v1
//     Version: 1.0.0
//     License: MIT https://raw.githubusercontent.com/e154/smart-home/master/LICENSE
//     Contact: Alex Filippov <support@e154.ru> https://e154.github.io/smart-home/
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - BasicAuth
//     - ApiKeyAuth
//
//     SecurityDefinitions:
//     ApiKeyAuth:
//          type: apiKey
//          name: Authorization
//          in: header
//     BasicAuth:
//          type: basic
//
// swagger:meta
package v1

import (
	_ "github.com/e154/smart-home/api/server/v1/controllers"
	_ "github.com/e154/smart-home/api/server/v1/models"
	_ "github.com/e154/smart-home/api/server/v1/responses"
)
