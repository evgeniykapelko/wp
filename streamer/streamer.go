package streamer

import (
	"fmt"
	"github.com/tsawler/toolbox"
	"path"
	"path/filepath"
	"strings"
)

type ProcessingMessage struct {
	ID         int
	Successful bool
	Message    string
	OutputFile string
}

type VideoProcessingJob struct {
	Video Video
}

type Processor struct {
	Engine Encoder
}

type Video struct {
	ID           int
	InputFile    string
	OutputDir    string
	EncodingType string
	NotifyChan   chan ProcessingMessage
	Options      *VideoOptions
	Encoder      Processor
}

type VideoOptions struct {
	RenameOutput    bool
	SegmentDuration int
	MaxRate1080p    string
	MaxRate720p     string
	MaxRate480p     string
}

func (vd *VideoDispatcher) NewVideo(id int, input, output, encType string, notifyChan chan ProcessingMessage, options *VideoOptions) Video {
	if options == nil {
		options = &VideoOptions{}
	}

	fmt.Println("New video created:", id, input)

	return Video{
		ID:           id,
		InputFile:    input,
		OutputDir:    output,
		EncodingType: encType,
		NotifyChan:   notifyChan,
		Encoder:      vd.Processor,
		Options:      options,
	}
}

func (v *Video) encode() {
	var fileName string

	switch v.EncodingType {
	case "mp4":
		fmt.Println("v.encode(): About to encode to mp4", v.ID)
		name, err := v.encodeToMP4()
		if err != nil {
			// send information to the notifyChan
			v.sendToNotifyChan(false, "", fmt.Sprintf("Error encoding to MP4: %s", err))
			return
		}
		fileName = fmt.Sprintf("%s.mp4", name)
	case "hls":
		name, err := v.encodeToHLS()
		if err != nil {
			// send information to the notifyChan
			v.sendToNotifyChan(false, "", fmt.Sprintf("Error encoding to MP4: %s", err))
			return
		}
		fileName = fmt.Sprintf("%s.m3u8", name)
	default:
		v.sendToNotifyChan(false, "", fmt.Sprintf("Encoding type not supported: %s", v.EncodingType))
		return
	}

	v.sendToNotifyChan(true, fileName, fmt.Sprintf("Encoding type %s", v.EncodingType))
}

func (v *Video) encodeToMP4() (string, error) {
	baseFileName := ""

	if !v.Options.RenameOutput {
		b := path.Base(v.InputFile)
		baseFileName = strings.TrimSuffix(b, filepath.Ext(b))
	} else {
		var t toolbox.Tools
		baseFileName = t.RandomString(10)

	}

	err := v.Encoder.Engine.EncodeToMP4(v, baseFileName)
	if err != nil {
		return "", err
	}

	return baseFileName, nil
}

func (v *Video) encodeToHLS() (string, error) {
	baseFileName := ""

	if !v.Options.RenameOutput {
		b := path.Base(v.InputFile)
		baseFileName = strings.TrimSuffix(b, filepath.Ext(b))
	} else {
		var t toolbox.Tools
		baseFileName = t.RandomString(10)
	}

	err := v.Encoder.Engine.EncodeToHLS(v, baseFileName)
	if err != nil {
		return "", err
	}

	return baseFileName, nil
}

func (v *Video) sendToNotifyChan(successful bool, fileName string, message string) {
	v.NotifyChan <- ProcessingMessage{
		ID:         v.ID,
		Successful: successful,
		Message:    message,
		OutputFile: fileName,
	}
}

func New(jobQueue chan VideoProcessingJob, maxWorkers int) *VideoDispatcher {
	workerPool := make(chan chan VideoProcessingJob, maxWorkers)

	var e VideoEncoder
	p := Processor{
		Engine: &e,
	}

	return &VideoDispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		WorkerPool: workerPool,
		Processor:  p,
	}
}
