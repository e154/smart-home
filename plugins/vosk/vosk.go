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

//go:build !test
// +build !test

package vosk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	vosk "github.com/alphacep/vosk-api/go"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/web"
)

type Vosk struct {
	modelPath      string
	modelName      string
	modelLoaded    *atomic.Bool
	model          *vosk.VoskModel
	recsmu         sync.Mutex
	grmRecs        *Pool
	gpRecs         *Pool
	crawler        web.Crawler
	isDownload     bool
	downloadCancel context.CancelFunc
}

func NewVosk(modelPath, modelName string, crawler web.Crawler) *Vosk {
	return &Vosk{
		modelPath:   modelPath,
		modelName:   modelName,
		modelLoaded: atomic.NewBool(false),
		crawler:     crawler,
	}
}

func (v *Vosk) Start() {

	go func() {
		if err := v.start(); err != nil {
			log.Error(err.Error())
		}
	}()
}

func (v *Vosk) start() error {

	if err := v.CheckModel(); err != nil {
		return err
	}

	if err := v.LoadModel(); err != nil {
		return err
	}

	v.grmRecs = NewPool(v, 1)
	v.gpRecs = NewPool(v, 1)

	if err := v.runTest(); err != nil {
		return err
	}

	log.Info("Vosk initiated successfully")

	return nil
}

func (v *Vosk) Shutdown() {
	if v.isDownload {
		v.downloadCancel()
	}
	if v.grmRecs != nil {
		v.grmRecs.Free()
	}
	if v.gpRecs != nil {
		v.gpRecs.Free()
	}
	if v.model != nil {
		v.model.Free()
	}
	v.modelLoaded.Store(false)
}

func (v *Vosk) CheckModel() error {
	modelName := v.prepareModelName(v.modelName)

	mPath := filepath.Join(v.modelPath, modelName)
	if _, err := os.Stat(mPath); err != nil {
		log.Warn("Path does not exist: " + mPath)
		if err = v.DownloadModel(modelName); err != nil {
			return fmt.Errorf("with model %s, error: %s", modelName, err.Error())
		}
	}
	return nil
}

func (v *Vosk) LoadModel() error {
	if v.modelLoaded.Load() {
		log.Info("A model was already loaded, freeing all recognizers and model")
		v.model.Free()
		v.modelLoaded.Store(false)
	}

	modelName := v.prepareModelName(v.modelName)
	mPath := filepath.Join(v.modelPath, modelName)
	log.Infof("Opening VOSK model (%s)", mPath)

	vosk.SetLogLevel(-1)
	aModel, err := vosk.NewModel(mPath)
	if err != nil {
		return err
	}
	v.model = aModel

	v.modelLoaded.Store(true)

	return nil
}

func (v *Vosk) STT(reader io.Reader, withGrm bool) (string, error) {

	worker := v.grmRecs.GetWorker()
	defer func() {
		v.grmRecs.Release(worker)
	}()

	buf := make([]byte, 4096)

	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Error(err.Error())
			}

			break
		}

		worker.Rec.AcceptWaveform(buf)
	}

	var jres map[string]interface{}
	err := json.Unmarshal(worker.Rec.Result(), &jres)
	if err != nil {
		return "", err
	}

	return jres["text"].(string), nil
}

func (v *Vosk) runTest() error {

	log.Info("Running recognizer test")

	sttTestPath := filepath.Join("data", "vosk", "stttest.pcm")
	pcmBytes, err := os.ReadFile(sttTestPath)
	if err != nil {
		return err
	}

	ioReaderData := bytes.NewReader(pcmBytes)
	buf := bytes.NewBuffer(make([]byte, 0))
	if _, err = buf.ReadFrom(ioReaderData); err != nil {
		return err
	}

	cTime := time.Now()

	transcribedText, err := v.STT(buf, false)
	if err != nil {
		return err
	}

	tTime := time.Now().Sub(cTime)
	log.Infof("Text (from test): %s, ok: %t", transcribedText, transcribedText == "how are you")
	if tTime.Seconds() > 3 {
		log.Infof("Vosk test took a while, performance may be degraded. (%s)", tTime)
		return nil
	}

	log.Infof("Vosk test successful! (Took %s)", tTime)

	return nil
}

func (v *Vosk) DownloadModel(modelName string) error {

	destpath := filepath.Join(modelPath)
	_ = os.MkdirAll(destpath, 0755)

	var ctx context.Context
	ctx, v.downloadCancel = context.WithCancel(context.Background())
	v.isDownload = true
	defer func() {
		v.isDownload = false
		v.downloadCancel = nil
	}()

	downloadedFilePath, err := v.crawler.Download(web.Request{
		Method:  "GET",
		Url:     URLPrefix + modelName + ".zip",
		Context: ctx,
	})
	if err != nil {
		return err
	}

	err = common.Unzip(downloadedFilePath, modelPath)
	if err != nil {
		return err
	}

	_ = os.Remove(downloadedFilePath)

	log.Info("Reloaded voice processor successfully")

	return nil
}

func (v *Vosk) prepareModelName(modelName string) string {
	modelName = strings.ReplaceAll(modelName, URLPrefix, "")
	modelName = strings.ReplaceAll(modelName, "https://", "")
	modelName = strings.TrimSuffix(modelName, ".zip")
	modelName = strings.ReplaceAll(modelName, "/", "")
	modelName = strings.TrimSpace(modelName)
	return modelName
}
