package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
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

func readFile(path string, prduct string) {
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
}

func main() {
	var home string = os.Getenv("HOME")
	flag.Parse()
	if appName != "" {
		fmt.Println(appName)
	} else {
		log.Fatal("Product Name not provided")
	}
	cb := make(callBackChan)
	go triggerReadFile(10*time.Second, cb)
	go func() {
		for {
			select {
			case <-cb:
				readFile(home, appName)
			}
		}
	}()
	for {
		time.Sleep(10 * time.Second)
	}
}
