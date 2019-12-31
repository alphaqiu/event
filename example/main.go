package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	event "github.com/alphaqiu/event"
)

const (
	accepted event.Type = "Accepted"
	lost     event.Type = "Lost"
)

func main() {
	ee := event.NewEventEmitter(7)
	go producer(ee)
	go consumer(1, ee)
	go consumer(2, ee)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}

func producer(ee *event.Emitter) {
	time.Sleep(time.Second * 2)
	fmt.Println("producer starting")
	for i := 0; i < 10; i++ {
		ee.Emit(event.NewEvent(accepted, i+1))
		fmt.Println("producer >>>", i+1)
	}
	fmt.Println("producer done...")
}

func consumer(index int, ee *event.Emitter) {
	fmt.Printf("consumer[%d] starting\n", index)
	var cancel event.Cancel
	counter := 0
	if index <= 1 {
		cancel = ee.On(func(e *event.Event) {
			// time.Sleep(time.Second)
			fmt.Printf("consumer[%d] fireAlways recv: %d\n", index, e.GetValue())
			counter += 1
			if counter > 5 {
				fmt.Println("cancel before")
				cancel()
				fmt.Println("cancel after")
			}
		}, accepted, lost)
	} else {
		cancel = ee.Once(func(e *event.Event) {
			fmt.Printf("consumer[%d] fireOnce recv: %d\n", index, e.GetValue())
		}, accepted)
	}

}
