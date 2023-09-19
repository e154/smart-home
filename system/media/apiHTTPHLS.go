package media

import (
	"bytes"
	"time"

	"github.com/deepch/vdk/format/ts"
	"github.com/gin-gonic/gin"
)

// HTTPAPIServerStreamHLSM3U8 send client m3u8 play list
func HTTPAPIServerStreamHLSM3U8(c *gin.Context) {


	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		log.Error(ErrorStreamNotFound.Error())
		return
	}

	if !RemoteAuthorization("HLS", c.Param("uuid"), c.Param("channel"), c.Query("token"), c.ClientIP()) {
		log.Error(ErrorStreamUnauthorized.Error())
		return
	}

	c.Header("Content-Type", "application/x-mpegURL")
	Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
	//If stream mode on_demand need wait ready segment's
	for i := 0; i < 40; i++ {
		index, seq, err := Storage.StreamHLSm3u8(c.Param("uuid"), c.Param("channel"))
		if err != nil {
			c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
			log.Error(err.Error())
			return
		}
		if seq >= 6 {
			_, err := c.Writer.Write([]byte(index))
			if err != nil {
				c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
				log.Error(err.Error())
				return
			}
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// HTTPAPIServerStreamHLSTS send client ts segment
func HTTPAPIServerStreamHLSTS(c *gin.Context) {

	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		log.Error(ErrorStreamNotFound.Error())
		return
	}
	codecs, err := Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.Error(err.Error())
		return
	}
	outfile := bytes.NewBuffer([]byte{})
	Muxer := ts.NewMuxer(outfile)
	Muxer.PaddingToMakeCounterCont = true
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.Error(err.Error())
		return
	}
	seqData, err := Storage.StreamHLSTS(c.Param("uuid"), c.Param("channel"), stringToInt(c.Param("seq")))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.Error(err.Error())
		return
	}
	if len(seqData) == 0 {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotHLSSegments.Error()})
		log.Error(ErrorStreamNotHLSSegments.Error())
		return
	}
	for _, v := range seqData {
		v.CompositionTime = 1
		err = Muxer.WritePacket(*v)
		if err != nil {
			c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
			log.Error(err.Error())
			return
		}
	}
	err = Muxer.WriteTrailer()
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.Error(err.Error())
		return
	}
	_, err = c.Writer.Write(outfile.Bytes())
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		log.Error(err.Error())
		return
	}

}
