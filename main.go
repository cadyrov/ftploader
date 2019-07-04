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
	http.ListenAndServe(":"+port, r)
}

func download(w http.ResponseWriter, r *http.Request) {
	if pth, ok := r.URL.Query()["path"]; ok && len(pth) > 0 {

		data, err := ioutil.ReadFile(pth[0])
		if err != nil {
			log.Println(err)
			return
		}
		type slicefls struct {
			Files []string `json:"files"`
		}
		sfls := slicefls{}
		errs := yaml.Unmarshal([]byte(data), &sfls)
		if errs != nil {
			log.Println(errs)
			return
		}
		copyFiles(sfls.Files)
	} else {
		w.Write([]byte("send path parameter"))
		w.WriteHeader(http.StatusBadRequest)
	}

}

func copyFiles(slicefls []string) {
	config := app.NewConfig()
	var wg sync.WaitGroup
	if slicefls != nil {
		for _, value := range slicefls {
			wg.Add(1)
			name := value
			go func() {
				oneFile(name, name, *config)
				fmt.Printf("file copied %v \n", name)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	fmt.Println("process finished")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func oneFile(fileIn string, fileOut string, cnf app.Config) (result string) {
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
	err = client.Retrieve(cnf.FTP.SourcePath+"/"+fileIn, buf)
	check(err)

	fl, err := os.Create(cnf.FTP.DestinationPath + "/" + fileOut)
	check(err)
	defer fl.Close()

	n2, err := fl.Write(buf.Bytes())
	check(err)
	result = fmt.Sprintf("wrote %v bytes in %v", n2, fileOut)
	return
}
