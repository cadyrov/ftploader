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

var logpath string = ""

func main() {
	port := "3000"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	r := mux.NewRouter()
	r.HandleFunc("/download", download).Methods("GET")
	fmt.Println("Listening...")
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func toLog(msg string) {
	config := app.NewConfig()
	if logpath == "" && config.FTP.LogPath != "" {
		logpath = config.FTP.LogPath
	}
	if logpath != "" {
		f, err := os.OpenFile(logpath,
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
		if lg, ok := r.URL.Query()["log"]; ok && len(lg) > 0 {
			logpath = lg[0]
		}
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
		copyFiles(sfls.Files)
		w.Write([]byte("Process finished"))
	} else {
		w.Write([]byte("send path parameter"))
		w.WriteHeader(http.StatusBadRequest)
	}

}

func copyFiles(slicefls []string) {
	toLog("process started")
	toLog("logpath:" + logpath)
	config := app.NewConfig()
	var wg sync.WaitGroup
	goroutines := make(chan int, config.FTP.Connections)
	// Read data from input channel
	countError := 0
	if slicefls != nil {
		toLog(fmt.Sprintf("%v files", len(slicefls)))
		for _, value := range slicefls {
			name := value
			goroutines <- 1

			wg.Add(1)
			go func(name string, goroutines <-chan int, countError *int, wg *sync.WaitGroup) {
				msg, success := oneFile(name, name, *config)
				if msg != "" {
					toLog(msg)
					if success == false {
						*countError++
					}
				}
				<-goroutines
				wg.Done()
			}(name, goroutines, &countError, &wg)
		}
	}
	wg.Wait()
	close(goroutines)
	toLog("process finished")
	toLog(fmt.Sprintf("ошибки: %v", countError))
}

func check(e error) {
	if e != nil {
		toLog(e.Error())
	}
}

func oneFile(fileIn string, fileOut string, cnf app.Config) (result string, success bool) {
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
		if err != nil {
			result = err.Error()
			return
		}

		// download to a buffer instead of file
		buf := new(bytes.Buffer)
		err = client.Retrieve(src, buf)
		if err != nil {
			result = err.Error()
			return
		}

		fl, err := os.Create(destination)
		if err != nil {
			result = err.Error()
			return
		}
		defer fl.Close()

		n2, err := fl.Write(buf.Bytes())
		if err != nil {
			result = err.Error()
			return
		}
		result = fmt.Sprintf("загрузка %v bytes in %v", n2, fileOut)
	}
	success = true
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
