package cameraClient

import (
	"crypto/tls"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"time"

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

	// connect to the server
	err = gc.Start()
	if err != nil {
		panic(err)
	}
	defer gc.Close()

	// find available medias
	desc, _, err := gc.Describe(u)
	if err != nil {
		panic(err)
	}

	// find the H264 media and format
	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		panic("media not found")
	}

	// setup RTP -> H264 decoder
	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		panic(err)
	}

	// setup H264 -> RGBA decoder
	h264Dec := &h264Decoder{}
	err = h264Dec.initialize()
	if err != nil {
		panic(err)
	}
	defer h264Dec.close()

	// if SPS and PPS are present into the SDP, send them to the decoder
	if forma.SPS != nil {
		h264Dec.decode([][]byte{forma.SPS})
	}
	if forma.PPS != nil {
		h264Dec.decode([][]byte{forma.PPS})
	}

	// setup a single media
	_, err = gc.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		panic(err)
	}

	firstRandomAccess := false
	saveCount := 0

	// called when a RTP packet arrives
	gc.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// extract access units from RTP packets
		au, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Printf("ERR: %v", err)
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
			panic(err)
		}

		// check for frame presence
		if img == nil {
			log.Printf("ERR: frame cannot be decoded")
			return
		}

		// convert frame to JPEG and save to file
		err = saveToFile(img)
		if err != nil {
			panic(err)
		}

		saveCount++
		if saveCount == 3 {
			log.Printf("saved 3 images, exiting")
			os.Exit(1)
		}
	})

	// start playing
	_, err = gc.Play(nil)
	if err != nil {
		panic(err)
	}

	// wait until a fatal error
	panic(gc.Wait())
}

func saveToFile(img image.Image) error {
	// create file
	fname := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10) + ".jpg"
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.Println("saving", fname)

	// convert to jpeg
	return jpeg.Encode(f, img, &jpeg.Options{
		Quality: 60,
	})
}
