/*
   http://stackoverflow.com/questions/20041392/how-does-makechan-bool-behave-differently-from-makechan-bool-1
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	intInputChan := make(chan int)
	done := make(chan bool)
	quit := make(chan bool)

	go worker(intInputChan, done, quit)

	for i := 1; i < 10; i++ {
		intInputChan <- i
	}

	close(intInputChan)
	<-done

	fmt.Println("Existing Main App... ")

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	<-sigChan

}

/*
	Send and receive on nil channels block
	Select never selects a blocking case.

	Refer to the video from Sameer Ajmani.
*/

func worker(input chan int, done chan bool, quit chan bool) {
	for {

		fmt.Println(".......Select........")

		/*
			Below select executes when an input is received on a channel
			from a source..
		*/

		select {
		case intVal, ok := <-input:
			if ok {
				fmt.Printf("--> %v %v\n", intVal, ok)
				time.Sleep(100 * time.Millisecond)
			} else {
				fmt.Println(" Input Channel closed..")

				/* Disable( turn off) input channel */
				input = nil

				time.Sleep(100 * time.Millisecond)
				/*
				 Standard way to close channel and send signal on the select
				 for correcponding channel and this will trigger <-quit below
				*/
				close(quit)

				/*
					This will not work in the context of current select
					because Select selects one channel input at a time
					and sending on that will frooze the select.
				*/
				/*---------------------------------------------------------*/
				/*
					For this to work make a buffered channel like quit = make(chan bool, 1)
					refer to the link at the top for info... inn anut shell..
					unbuffered channel is block moment we pass in a data to it till somesone
					read data for that whereas buffered does not get blocked and is available
					for the below case to read and respond
				*/
				//quit <- true
			}
		case b, ok := <-quit:
			if ok {
				fmt.Println("quit ok normal", b)
				//done <- true
				//return
			} else {
				fmt.Println(`"quit ok" on channel closed..`, b)

				/* Disable( turn off) quit channel */
				quit = nil
				done <- true

				/* return takes us out of the indefinite for loop*/
				return
			}
		}
	}

	fmt.Println(".....Out of for loop....")
}
