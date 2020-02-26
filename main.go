package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var changeChan = make(chan string)

// @TODO directories watcher and file-watcher update
// @TODO directory excludes

func main() {
	var location string
	if len(os.Args) != 2 {
		// log.Fatal("You need set location as a first argument")
		// dev
		location = "src"
	} else {
		location = os.Args[1]
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	http.HandleFunc("/ws", serveServiceWorker)

	// dev
	//fs := http.FileServer(http.Dir("src"))
	//http.Handle("/", fs)
	// dev

	fmt.Println("Watching:", location)
	files, _ := getWatchedFiles(location)
	go runWatcher(files, location)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveServiceWorker(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	log.Println("Connection from:", r.RemoteAddr)
	go ruotineWrite(ws)
}

func runWatcher(files []string, location string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("modified OP:", event.Op)
				if (event.Op & fsnotify.Write) == fsnotify.Write {
					log.Println("modified file:", event.Name)
					// @TODO trottling
					changeChan <- event.Name[len(location)+1:]
				}
				// @TODO try to restore file when -> rename -> chmod -> write
				// for vim

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done
}

func getWatchedFiles(root string) (files []string, directories []string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			directories = append(directories, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return
}

func ruotineWrite(ws *websocket.Conn) {
	defer func() {
		ws.Close()
		log.Println(ws.RemoteAddr())
	}()
	for {
		select {
		case name := <-changeChan:
			log.Println("WebSocket message:", name)
			ws.SetWriteDeadline(time.Now().Add(1 * time.Second))
			err := ws.WriteMessage(websocket.TextMessage, []byte(name))
			if err != nil {
				log.Println("error:", err)
			}
		}
	}
}
