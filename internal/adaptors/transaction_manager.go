// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package adaptors

import (
	"context"
	"fmt"

	"github.com/e154/smart-home/internal/db"
	"gorm.io/gorm"
)

func InjectTransaction(ctx context.Context, tr *gorm.DB) context.Context {
	return context.WithValue(ctx, db.GormTransaction, tr)
}

type TransactionManger struct {
	db *gorm.DB
}

func NewTransactionManger(db *gorm.DB) *TransactionManger {
	return &TransactionManger{db: db}
}

func (m *TransactionManger) Do(ctx context.Context, fn func(context.Context) error) (doErr error) {
	tr := m.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tr.Rollback()
			if doErr == nil {
				doErr = fmt.Errorf("Recover")
			}
			return
		}

		if doErr != nil {
			tr.Rollback()
			return
		}
		tr.Commit()
	}()

	return fn(InjectTransaction(ctx, tr))
}
