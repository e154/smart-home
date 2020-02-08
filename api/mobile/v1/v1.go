// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

// Smart home
//
// This documentation describes APIs found under https://github.com/e154/smart-home
//
//     BasePath: /api/v1
//     Version: 1.0.0
//     License: GPLv3 https://raw.githubusercontent.com/e154/smart-home/master/LICENSE
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
	_ "github.com/e154/smart-home/api/mobile/v1/controllers"
	_ "github.com/e154/smart-home/api/server/v1/models"
	_ "github.com/e154/smart-home/api/server/v1/responses"
)
