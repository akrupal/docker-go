package main

import "log"

func main() {
	log.Println("Hello World!")
	// to build image last part is optional but is better than the random tag
	// docker build . -t docker-containerised:latest
	// to list all the images
	// docker image ls
	// to run the code in docker image
	// docker run docker-containerised:latest
}
