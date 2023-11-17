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
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func splitFile(inputPath string, chunkSize int64) ([]string, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	chunkSize = 1024 * 1024 * chunkSize

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	numChunks := (fileSize + chunkSize - 1) / chunkSize

	tmpDir := path.Join(os.TempDir(), "smart_home")
	if err = os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, err
	}

	var fileList []string
	fileName := strings.ReplaceAll(fileInfo.Name(), ".zip", "")
	for i := int64(0); i < numChunks; i++ {
		outputPath := fmt.Sprintf("%s_part%d.dat", fileName, i+1)
		outputPath = path.Join(tmpDir, outputPath)

		// Создаем новый файл для части
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return nil, err
		}
		defer outputFile.Close()

		// Читаем кусок данных из исходного файла
		bufferSize := chunkSize
		if i == numChunks-1 {
			bufferSize = fileSize - i*chunkSize
		}

		buffer := make([]byte, bufferSize)
		_, err = inputFile.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		// Записываем кусок в новый файл
		_, err = outputFile.Write(buffer)
		if err != nil {
			return nil, err
		}

		fileList = append(fileList, outputPath)
	}

	return fileList, nil
}

func joinFiles(inputPattern, dir string) (string, error) {

	matches, err := filepath.Glob(inputPattern)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 || len(matches) == 1 {
		return "", errors.New("no matches or less then 2")
	}

	sort.Strings(matches)

	fmt.Println(matches)

	params := strings.Split(matches[0], "_part")
	if len(params) < 2 {
		return "", err
	}

	fileName := path.Join(dir, filepath.Base(fmt.Sprintf("%s.zip", params[0])))
	outputFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	for _, match := range matches {
		inputFile, err := os.Open(match)
		if err != nil {
			return "", err
		}
		defer inputFile.Close()

		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			return "", err
		}
	}

	return fileName, nil
}
