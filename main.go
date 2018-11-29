package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

//const maxUploadSize = 8 * 1024 * 1024 // 2 mb
const uploadPath = "./tmp"

func main() {
	s := NewStorage(uploadPath)
	r := chi.NewRouter() //r.Use(middleware.RequestID)
	r.Use(CrossControl)
	fs := FileServer{s: s}
	http.HandleFunc("/upload", fs.UploadFile)
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir(uploadPath))))
	//curl -i  -F 'uploadfile=@/home/egor/Work/Golang/src/github.com/Jopoleon/RuslanTest/1.jpg' http://localhost:8080/upload
	logrus.Info("Server started on localhost:8080, use /upload for uploading files and /files/{fileName} for downloading")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CrossControl(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func uploadFileHandler2() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		f, err := os.Create("./tmp/" + handler.Filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		//f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		//if strings.Contains(err.Error(),"no such file or directory") {
		//	err := os.Mkdir("./test/"+handler.Filename, 0666)
		//	if err != nil {
		//		fmt.Println(err)
		//		return
		//	}
		//} else {
		//	fmt.Println(err)
		//	return
		//}
		io.Copy(f, file)
		w.Write([]byte("SUCCESS"))
	})

}

//
//func uploadFileHandler() http.HandlerFunc {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// validate file size
//
//		logrus.Info(r)
//		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
//		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
//			renderError(w, "FILE_TOO_BIG \n", http.StatusBadRequest)
//			return
//		}
//
//		// parse and validate file and post parameters
//		fileType := r.PostFormValue("type")
//		file, _, err := r.FormFile("uploadFile")
//		if err != nil {
//			logrus.Error(err)
//			renderError(w, "INVALID_FILE \n", http.StatusBadRequest)
//			return
//		}
//		defer file.Close()
//		fileBytes, err := ioutil.ReadAll(file)
//		if err != nil {
//			logrus.Error(err)
//			renderError(w, "INVALID_FILE \n", http.StatusBadRequest)
//			return
//		}
//
//		// check file type, detectcontenttype only needs the first 512 bytes
//		filetype := http.DetectContentType(fileBytes)
//		switch filetype {
//		case "image/jpeg", "image/jpg":
//		case "image/gif", "image/png":
//		case "application/pdf":
//			break
//		default:
//			renderError(w, "INVALID_FILE_TYPE \n", http.StatusBadRequest)
//			return
//		}
//		fileName := randToken(12)
//		fileEndings, err := mime.ExtensionsByType(fileType)
//		if err != nil {
//			renderError(w, "CANT_READ_FILE_TYPE \n", http.StatusInternalServerError)
//			return
//		}
//		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
//		fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)
//
//		// write file
//		newFile, err := os.Create(newPath)
//		if err != nil {
//			renderError(w, "CANT_WRITE_FILE \n", http.StatusInternalServerError)
//			return
//		}
//		defer newFile.Close() // idempotent, okay to call twice
//		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
//			renderError(w, "CANT_WRITE_FILE \n", http.StatusInternalServerError)
//			return
//		}
//		w.Write([]byte("SUCCESS"))
//	})
//}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
