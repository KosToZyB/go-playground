package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/satori/go.uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
var storageDir = "./image-storage/"

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/", http.FileServer(http.Dir(storageDir)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	u, err := uuid.NewV4()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fileName := u.String()

	pathToUpload := getPathToUpload(fileName)

	fileName += path.Ext(header.Filename)

	if _, err := os.Stat(storageDir + pathToUpload); os.IsNotExist(err) {

		err = os.MkdirAll(storageDir+pathToUpload, 0755)
		if err != nil {
			fmt.Fprintf(w, "Error create dir")
			return
		}
	}

	filePath := storageDir + pathToUpload + fileName
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully. URL: ")
	fmt.Fprintf(w, pathToUpload+fileName)
}

func getPathToUpload(fileName string) string {
	const levelNested = 2
	const lenDirName = 2
	const sep = "/"

	var filePath string
	if len(fileName) < levelNested*lenDirName {
		for i := 0; i < levelNested; i++ {
			filePath += randStringRunes(lenDirName)
			filePath += sep
		}
		return filePath
	}

	for i := 0; i < levelNested; i++ {
		filePath += fileName[i*lenDirName : (i+1)*lenDirName]
		filePath += sep
	}

	return filePath
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
