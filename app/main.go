package main

import (
	"fmt"
	"streamer"
)

func main() {
	const numJobs = 4
	const numWorkers = 2

	notifyChan := make(chan streamer.ProcessingMessage, numJobs)
	defer close(notifyChan)

	videoQueue := make(chan streamer.VideoProcessingJob, numJobs)
	defer close(videoQueue)

	wp := streamer.New(videoQueue, numWorkers)
	fmt.Println("wp:", wp)

	wp.Run()
}
