// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package backup

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitFile(t *testing.T) {

	inputPath := "../../snapshots/2023-11-06T15:12:55.284.zip"
	var chunkSize int64 = 25 // 25 MB

	fileList, err := splitFile(inputPath, chunkSize)
	require.NoError(t, err)
	require.NotNil(t, fileList)
	fmt.Println(fileList)
}


func TestJoinFiles(t *testing.T) {

	inputPattern := "*_part*.dat"

	fileName, err := joinFiles(inputPattern)
	require.NoError(t, err)
	fmt.Println(fileName)
}


func TestCheckZip(t *testing.T) {

	inputPattern := "2023-11-06T15:12:55.284.zip"

	_, err := checkZip(inputPattern)
	require.NoError(t, err)
}
