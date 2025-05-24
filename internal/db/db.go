// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package db

import (
	"context"

	"gorm.io/gorm"
)

type TransactionKey string

const GormTransaction TransactionKey = "gorm_transaction"

func ExtractTransaction(ctx context.Context) *gorm.DB {
	tr, ok := ctx.Value(GormTransaction).(*gorm.DB)
	if !ok {
		return nil
	}

	return tr
}

type Common struct {
	Db *gorm.DB
}

func (c *Common) DB(ctx context.Context) *gorm.DB {

	db := c.Db
	if tr := ExtractTransaction(ctx); tr != nil {
		db = tr
	}

	return db.WithContext(ctx)
}
