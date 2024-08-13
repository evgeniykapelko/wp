package main

import (
	"fmt"
	"streamer"
)

func main() {
	const numJobs = 1
	const numWorkers = 2

	notifyChan := make(chan streamer.ProcessingMessage, numJobs)
	defer close(notifyChan)

	videoQueue := make(chan streamer.VideoProcessingJob, numJobs)
	defer close(videoQueue)

	wp := streamer.New(videoQueue, numWorkers)

	wp.Run()
	fmt.Println("Worker pool started. Press enter to continue.")
	_, _ = fmt.Scanln()

	video := wp.NewVideo(1, "./input/puppy1.mp4", "./output", "mp4", notifyChan, nil)

	videoQueue <- streamer.VideoProcessingJob{Video: video}

	for i := 1; i < numJobs; i++ {
		msg := <-notifyChan
		fmt.Println(msg)
	}

	fmt.Println("Done!")
}
