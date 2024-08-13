package main

import (
	"fmt"
	"streamer"
)

func main() {
	const numJobs = 4
	const numWorkers = 4

	notifyChan := make(chan streamer.ProcessingMessage, numJobs)
	defer close(notifyChan)

	videoQueue := make(chan streamer.VideoProcessingJob, numJobs)
	defer close(videoQueue)

	wp := streamer.New(videoQueue, numWorkers)

	wp.Run()

	video1 := wp.NewVideo(1, "./input/puppy1.mp4", "./output", "mp4", notifyChan, nil)

	video2 := wp.NewVideo(2, "./input/bad.txt", "./output", "mp4", notifyChan, nil)
	ops := &streamer.VideoOptions{
		RenameOutput:    true,
		SegmentDuration: 10,
		MaxRate1080p:    "1200k",
		MaxRate720p:     "600K",
		MaxRate480p:     "400K",
	}
	video3 := wp.NewVideo(3, "./input/puppy2.mp4", "./output", "hls", notifyChan, ops)

	video4 := wp.NewVideo(4, "./input/puppy2.mp4", "./output", "mp4", notifyChan, nil)

	videoQueue <- streamer.VideoProcessingJob{Video: video1}
	videoQueue <- streamer.VideoProcessingJob{Video: video2}
	videoQueue <- streamer.VideoProcessingJob{Video: video3}
	videoQueue <- streamer.VideoProcessingJob{Video: video4}

	for i := 1; i < numJobs; i++ {
		msg := <-notifyChan
		fmt.Println(msg)
	}

	fmt.Println("Done!")
}
