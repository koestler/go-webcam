package cameraClient

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type rtspState struct {
	tmpDir string
}

func createRtspState() rtspState {
	tmpDir, err := os.MkdirTemp("", "go-webcam")
	if err != nil {
		log.Fatalf("failed to create rtspState tmpDir: %v", err)
	}

	return rtspState{
		tmpDir: tmpDir,
	}
}

func (rs *rtspState) Close() {
	err := os.RemoveAll(rs.tmpDir)
	if err != nil {
		log.Printf("failed to remove rtspState tmpDir=%s: %v", rs.tmpDir, err)
	}
}

func (c *Client) rtspGetRawImage() (img []byte, err error) {
	url := c.config.Address()
	outputFile := filepath.Join(c.rtsp.tmpDir, "tmp.jpg")

	cmd := exec.Command("ffmpeg",
		"-y",
		"-threads", "1",
		"-rtsp_transport", "tcp",
		"-i", url,
		"-vframes", "1",
		outputFile,
	)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return os.ReadFile(outputFile)
}
