//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream) <-chan *Tweet {
	tw := make(chan *Tweet)
	go func() {
		defer close(tw)
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				break
			}
			if tweet == nil {
				break
			}
			tw <- tweet
		}
		fmt.Println("closed")
	}()
	return tw
}

func consumer(tweet *Tweet) {
	if tweet.IsTalkingAboutGo() {
		fmt.Println(tweet.Username, "\ttweets about golang")
	} else {
		fmt.Println(tweet.Username, "\tdoes not tweet about golang")
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	tw := producer(stream)

	// Consumer
	for {
		select {
		case tweet, ok := <-tw:
			if ok {
				consumer(tweet)
			} else {
				fmt.Printf("Process took %s\n", time.Since(start))
				return
			}
		}
	}

}
