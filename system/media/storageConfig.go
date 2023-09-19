package media

import (
	"time"

	"github.com/imdario/mergo"
)

// Command line flag global variables
var configFile string

// NewStreamCore do load config file
func NewStreamCore() *StorageST {
	//configFile = path.Join("conf", "media_config.json")
	//
	var instance StorageST
	//data, err := ioutil.ReadFile(configFile)
	//if err != nil {
	//	log.Error(err.Error())
	//	os.Exit(1)
	//}
	//err = json.Unmarshal(data, &tmp)
	//if err != nil {
	//	log.Error(err.Error())
	//	os.Exit(1)
	//}
	//debug = tmp.Server.Debug
	instance.Server = ServerST{
		Debug:              false,
		HTTPDemo:           true,
		HTTPDebug:          false,
		HTTPLogin:          "demo",
		HTTPPassword:       "demo",
		HTTPDir:            "static_source/media",
		HTTPPort:           ":8083",
		RTSPPort:           ":5541",
	}
	instance.Streams = make(map[string]StreamST)
	var err error
	for i, stream := range instance.Streams {
		for j, ch := range stream.Channels {
			channel := instance.ChannelDefaults
			err = mergo.Merge(&channel, ch)
			if err != nil {
				log.Error(err.Error())
				//os.Exit(1)
			}
			channel.clients = make(map[string]ClientST)
			channel.ack = time.Now().Add(-255 * time.Hour)
			channel.hlsSegmentBuffer = make(map[int]SegmentOld)
			channel.signals = make(chan int, 100)
			stream.Channels[j] = channel
		}
		instance.Streams[i] = stream
	}
	return &instance
}

// SaveConfig ...
func (obj *StorageST) SaveConfig() error {
	//log.Debug("Saving configuration to", configFile)
	//debug.Println(obj)
	//v2, err := version.NewVersion("2.0.0")
	//if err != nil {
	//	return err
	//}
	//data, err := sheriff.Marshal(&sheriff.Options{
	//	Groups:     []string{"config"},
	//	ApiVersion: v2,
	//}, obj)
	//if err != nil {
	//	return err
	//}
	//res, err := json.MarshalIndent(data, "", "  ")
	//if err != nil {
	//	return err
	//}
	//err = ioutil.WriteFile(configFile, res, 0644)
	//if err != nil {
	//	log.Error(err.Error())
	//	return err
	//}
	return nil
}
