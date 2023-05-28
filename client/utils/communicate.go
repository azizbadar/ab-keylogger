package utils

import (
	"bufio"
	"bytes"
	"os"
)

// var (
// 	userDirectory, _  = os.UserHomeDir()
// 	userTempDirectory = os.TempDir()
// 	userArchi         = runtime.GOARCH
// 	userOs            = string(runtime.GOOS)
// 	hostname, _       = os.Hostname()
// )

// func SendDataToServer() {

// }

// func SysInfo() string {

// 	res, _ := http.Get("https://api.ipify.org")
// 	ip, _ := io.ReadAll(res.Body)
// 	sysInfo := ("Operating System:	" + userOs + "\nHostname:	" + hostname + "\nPublic Ip:	" + string(ip) + "\nArchitecture:	" + userArchi + "\nUser Home Directory:	" + userDirectory + "\nUser Temp Directory:	" + userTempDirectory + "\n------------------------------------\nKEYS\n------------------------------------\n")
// 	sys := base64.StdEncoding.EncodeToString([]byte(sysInfo))
// 	return sys
// }

func PreProcessImage(dat *os.File) (*bytes.Reader, error) {
	stats, err := dat.Stat()
	if err != nil {
		return nil, err
	}
	var size = stats.Size()
	b := make([]byte, size)
	bufR := bufio.NewReader(dat)
	_, err = bufR.Read(b)
	bReader := bytes.NewReader(b)
	return bReader, err
}
