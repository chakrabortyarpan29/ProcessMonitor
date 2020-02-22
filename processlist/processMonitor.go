package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type callbackchan chan struct{}

// triggers after each time duration
func checkEveryInterval(ctx context.Context, d time.Duration, cb callbackchan) {
	for {
		select {
		case <-ctx.Done():
			// ctx is canceled
			return
		case <-time.After(d):
			// wait for the duration
			if cb != nil {
				cb <- struct{}{}
			}
		}
	}

}

// PrintProcessList for checking the process on machine
func PrintProcessList() {
	psCommand := exec.Command("ps", "a")
	resp, err := psCommand.CombinedOutput()
	if err != nil {
		log.Fatal("error: ps command failed")
	}

	out := string(resp)
	lines := strings.Split(out, "\n")

	for _, line := range lines {
		if line != "" {
			fmt.Println(line)
		}
	}
}

func main() {
	ctx := context.Background()
	PrintProcessList()
	cb := make(callbackchan)
	go checkEveryInterval(ctx, 5*time.Second, cb)
	go func() {
		for {
			select {
			case <-cb:
				PrintProcessList()
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}

}
