package cameraClient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"

	"github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/bluenviron/gortsplib/v5/pkg/format/rtph264"
	"github.com/bluenviron/mediacommon/v2/pkg/codecs/h264"
	"github.com/pion/rtp"
)

type rtspState struct {
}

func createRtspState() rtspState {
	return rtspState{}
}

func (c *Client) getRawImage() (img []byte, err error) {
	u, err := base.ParseURL(c.Config().Address())
	if err != nil {
		return nil, err
	}

	gc := gortsplib.Client{
		Scheme: u.Scheme,
		Host:   u.Host,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // TODO: make this configurable
		},
	}

	err = gc.Start()
	if err != nil {
		return nil, err
	}
	defer gc.Close()

	desc, _, err := gc.Describe(u)
	if err != nil {
		return nil, err
	}

	log.Printf("desc: %v", desc)

	for _, medi := range desc.Medias {
		log.Printf("media: %v", medi)
		for _, forma := range medi.Formats {
			fmt.Printf("codec: %v\n", forma.Codec())
		}
	}

	// find the H264 media and format
	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		return nil, errors.New("media not found")
	}

	// setup RTP -> H264 decoder
	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		return nil, err
	}

	// setup H264 -> RGBA decoder
	h264Dec := &h264Decoder{}
	err = h264Dec.initialize()
	if err != nil {
		return nil, err
	}
	defer h264Dec.close()

	// if SPS and PPS are present into the SDP, send them to the decoder
	if forma.SPS != nil {
		_, err := h264Dec.decode([][]byte{forma.SPS})
		if err != nil {
			return nil, err
		}
	}
	if forma.PPS != nil {
		_, err := h264Dec.decode([][]byte{forma.PPS})
		if err != nil {
			return nil, err
		}
	}

	// set up a single media
	_, err = gc.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		return nil, err
	}

	firstRandomAccess := false

	// called when a RTP packet arrives
	gc.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// extract access units from RTP packets
		au, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Printf("rtpDec error: %v", err)
			}
			return
		}

		// wait for a random access unit
		if !firstRandomAccess && !h264.IsRandomAccess(au) {
			log.Printf("waiting for a random access unit")
			return
		}
		firstRandomAccess = true

		// convert H264 access units into RGBA frames
		img, err := h264Dec.decode(au)
		if err != nil {
			log.Printf("h264Dec error: %v", err)
			return
		}

		// check for frame presence
		if img == nil {
			log.Printf("ERR: frame cannot be decoded")
			return
		}
		// convert frame to JPEG and save to file
		//err = saveToFile(img)
		//if err != nil {
		//	panic(err)
		//}
		fmt.Printf("********* VALID FRAME DECODED **********\n")

	})

	// start playing
	_, err = gc.Play(nil)
	if err != nil {
		return nil, err
	}

	// wait until a fatal error
	return nil, gc.Wait()
}
