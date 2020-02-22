package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	appName string
	errFoo  = errors.New("APP FOLDER NOT FOUND")
)

type callBackChan chan struct{}

const defaultPath = "."

func init() {
	flag.StringVar(&appName, "p", "", "Product Name")
}

func triggerReadFile(d time.Duration, cb callBackChan) {
	for {
		select {
		case <-time.After(d):
			{
				cb <- struct{}{}
			}
		}
	}
}

func readFile(path string, prduct string, wg *sync.WaitGroup) {
	path = path + "/" + prduct + ".txt"
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("File Not Found")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	wg.Done()
}

func main() {
	var home string = os.Getenv("HOME")
	flag.Parse()
	if appName != "" {
		fmt.Println(appName)
	} else {
		log.Fatal("Product Name not provided")
	}
	appLsit := strings.Split(appName, ",")
	numberOfWork := len(appLsit)
	var wg sync.WaitGroup
	cb := make(callBackChan)
	go triggerReadFile(10*time.Second, cb)
	for {
		select {
		case <-cb:
			for i := 0; i < numberOfWork; i++ {
				wg.Add(1)
				go readFile(home, appLsit[i], &wg)
			}
		}
	}
	wg.Wait()
}
