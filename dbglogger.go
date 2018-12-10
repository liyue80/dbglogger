package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type logMessage struct {
	Severity int    `json:"severity"`
	Content  string `json:"content"`
}

type serverConfig struct {
	PrintConsole    bool
	ConsoleSeverity int
	PrintFile       bool
	FileName        string
	FileSeverity    int
}

var defaultConfig = serverConfig{
	PrintConsole:    true,
	ConsoleSeverity: 0,
	PrintFile:       false,
	FileName:        "",
	FileSeverity:    0}

var config serverConfig
var mutex = &sync.Mutex{}
var fileHandle *os.File

const programHelpMsg = `
Ex:
  dbglogger [-c <filename>]
  `

func postMembersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m logMessage
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &m)

	mutex.Lock()
	defer mutex.Unlock()

	prefixTimestamp := time.Now().Format("2006-01-02 15:04:05.000")

	if config.PrintConsole {
		if config.ConsoleSeverity == 0 || config.ConsoleSeverity >= m.Severity {
			fmt.Printf("[%s] [%d] %s\n", prefixTimestamp, m.Severity, m.Content)
		}
	}

	if fileHandle != nil {
		if config.FileSeverity == 0 || config.FileSeverity >= m.Severity {
			fileHandle.WriteString(fmt.Sprintf("[%s] [%d] %s\n", prefixTimestamp, m.Severity, m.Content))
		}
	}

	// Send response
	w.Write([]byte("{}"))
}

func loadConfig() {
	switch len(os.Args) {
	case 1:
		fmt.Println("Loading default setting..")
		config.PrintConsole = defaultConfig.PrintConsole
		config.ConsoleSeverity = defaultConfig.ConsoleSeverity
		config.PrintFile = defaultConfig.PrintFile
		config.FileName = defaultConfig.FileName
		config.FileSeverity = defaultConfig.FileSeverity
	case 3:
		if os.Args[1] == "-c" {
			file, err := os.Open(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			fileInfo, _ := file.Stat()
			data := make([]byte, fileInfo.Size())
			file.Read(data)
			fmt.Println(string(data))
			if err = json.Unmarshal(data, &config); err != nil {
				log.Fatal("Parse config file error: ", err)
			}
			data = nil
		} else {
			fmt.Printf("Unknown option '%s'\n", os.Args[0])
			fmt.Println(programHelpMsg)
			os.Exit(1)
		}
	default:
		fmt.Println("Options error")
		fmt.Println(programHelpMsg)
		os.Exit(1)
	}

	// Open the log file
	if config.PrintFile && len(config.FileName) > 0 {
		var err error
		fileHandle, err = os.OpenFile(config.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		// No close this fileHandle until end of process
		if err != nil {
			log.Fatal(err)
		}
		fileHandle.WriteString("====== Debug Logger Starts ======\n")
	}
}

func main() {
	loadConfig()

	r := mux.NewRouter()
	r.HandleFunc("/dbgloggers", postMembersHandler).Methods("POST")

	http.Handle("/", r)

	fmt.Println("Listening on port 27109")
	if err := http.ListenAndServe(":27109", nil); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
