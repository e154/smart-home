package media

import (
	"fmt"
	"os"
	"time"

	"github.com/deepch/vdk/format/mp4"
	"github.com/gin-gonic/gin"
)

// HTTPAPIServerStreamSaveToMP4 func
func HTTPAPIServerStreamSaveToMP4(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()

	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		log.Error(ErrorStreamNotFound.Error())
		return
	}

	if !RemoteAuthorization("save", c.Param("uuid"), c.Param("channel"), c.Query("token"), c.ClientIP()) {
		log.Error(ErrorStreamUnauthorized.Error())
		return
	}
	c.Writer.Write([]byte("await save started"))
	go func() {
		Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
		cid, ch, _, err := Storage.ClientAdd(c.Param("uuid"), c.Param("channel"), MSE)
		if err != nil {
			log.Error(err.Error())
			return
		}

		defer Storage.ClientDelete(c.Param("uuid"), cid, c.Param("channel"))
		codecs, err := Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = os.MkdirAll(fmt.Sprintf("save/%s/%s/", c.Param("uuid"), c.Param("channel")), 0755)
		if err != nil {
			log.Error(err.Error())
		}
		f, err := os.Create(fmt.Sprintf("save/%s/%s/%s.mp4", c.Param("uuid"), c.Param("channel"), time.Now().String()))
		if err != nil {
			log.Error(err.Error())
		}
		defer f.Close()

		muxer := mp4.NewMuxer(f)
		err = muxer.WriteHeader(codecs)
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer muxer.WriteTrailer()

		var videoStart bool
		controlExit := make(chan bool, 10)
		dur, err := time.ParseDuration(c.Param("duration"))
		if err != nil {
			log.Error(err.Error())
		}
		saveLimit := time.NewTimer(dur)
		noVideo := time.NewTimer(10 * time.Second)
		defer log.Info("client exit")
		for {
			select {
			case <-controlExit:
				log.Error("Saved Reader End")
				return
			case <-saveLimit.C:
				log.Error("Saved Limit End")
				return
			case <-noVideo.C:
				log.Error(ErrorStreamNoVideo.Error())
				return
			case pck := <-ch:
				if pck.IsKeyFrame {
					noVideo.Reset(10 * time.Second)
					videoStart = true
				}
				if !videoStart {
					continue
				}
				if err = muxer.WritePacket(*pck); err != nil {
					return
				}
			}
		}
	}()
}
