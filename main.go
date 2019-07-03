package main

import (
	"bytes"
	"fmt"
	"ftploader/app"
	"github.com/secsy/goftp"
	"os"
	"sync"
	"time"
)

func main(){
	slicefls := []string{"0075566005376510.png", "0075566005391664.png"}
	copyFiles(slicefls)
}

func copyFiles(slicefls []string) {
	var wg sync.WaitGroup
	for key, _ := range slicefls {
		wg.Add(1)
		go func() {
			oneFile(slicefls[key],slicefls[key])
			fmt.Printf("file copied %v \n", slicefls[key])
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("process finished")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func oneFile(fileIn string, fileOut string ) (result string) {
	cnf := app.Config()
	config := goftp.Config{
		User:               cnf.FTP.Login,
		Password:           cnf.FTP.Password,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	client, err := goftp.DialConfig(config, cnf.FTP.Host)
	check(err)

	// download to a buffer instead of file
	buf := new(bytes.Buffer)
	err = client.Retrieve(cnf.FTP.StartPath + "/" + fileIn, buf)
	check(err)

	fl, err := os.Create(cnf.FTP.DestinationPath + "/" + fileOut)
	check(err)
	defer fl.Close()

	n2, err := fl.Write(buf.Bytes())
	check(err)
	result = fmt.Sprintf("wrote %v bytes in %v", n2, fileOut)
	return
}
