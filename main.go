package main

import (
	"bytes"
	"fmt"
	"ftploader/app"
	"github.com/gorilla/mux"
	"github.com/secsy/goftp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	port := "3000"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	r := mux.NewRouter()
	r.HandleFunc("/download", download).Methods("GET")
	log.Println("Listening...")
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		toLog(err.Error())
	}
}

func toLog(msg string) {
	config := app.NewConfig()
	if config.FTP.LogPath != "" {
		f, err := os.OpenFile(config.FTP.LogPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		logger := log.New(f, "ftploader", log.LstdFlags)
		logger.Println(msg)
	} else {
		log.Println(msg)
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	if pth, ok := r.URL.Query()["path"]; ok && len(pth) > 0 {

		data, err := ioutil.ReadFile(pth[0])
		if err != nil {
			toLog(err.Error())
			return
		}
		type slicefls struct {
			Files []string `json:"files"`
		}
		sfls := slicefls{}
		errs := yaml.Unmarshal([]byte(data), &sfls)
		if errs != nil {
			toLog(errs.Error())
			return
		}
		toLog(fmt.Sprintf("%v files created", len(sfls.Files)))
		copyFiles(sfls.Files)
	} else {
		w.Write([]byte("send path parameter"))
		w.WriteHeader(http.StatusBadRequest)
	}

}

func copyFiles(slicefls []string) {
	toLog("process started")
	config := app.NewConfig()
	var wg sync.WaitGroup
	goroutines := make(chan int, 10)
	// Read data from input channel
	if slicefls != nil {
		for _, value := range slicefls {
			toLog(value)
			name := value
			goroutines <- 1
			wg.Add(1)
			go func(name string, goroutines <-chan int, wg *sync.WaitGroup) {
				msg := oneFile(name, name, *config)
				toLog(msg)
				<-goroutines
				wg.Done()
			}(name, goroutines, &wg)
		}
	}
	wg.Wait()
	close(goroutines)
	toLog("process finished")
}

func check(e error) {
	if e != nil {
		toLog(e.Error())
	}
}

func oneFile(fileIn string, fileOut string, cnf app.Config) (result string) {
	src := cnf.FTP.SourcePath + "/" + fileIn
	destination := cnf.FTP.DestinationPath + "/" + fileOut
	ok, err := exists(destination)
	isExist := false
	if ok && err == nil {
		isExist = true
	}
	if cnf.FTP.IsRewrite == true || !isExist {
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
		err = client.Retrieve(src, buf)
		check(err)

		fl, err := os.Create(destination)
		check(err)
		defer fl.Close()

		n2, err := fl.Write(buf.Bytes())
		check(err)
		result = fmt.Sprintf("wrote %v bytes in %v", n2, fileOut)
	} else {
		result = fmt.Sprintf("file %v exists work done", destination)
	}
	return
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
