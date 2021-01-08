package main

import (
	"bytes"
	"fmt"
	"go_aes_server/crypto"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	PORT = "8080"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile(Pwd() + "html/index.html")
		if err != nil {
			log.Println("ioutil.ReadFile(Pwd() + \"html/index.html\") Exception:" + err.Error())
		}
		w.Write(bytes)
		return
	})

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile(Pwd() + "favicon.ico")
		if err != nil {
			log.Println("ioutil.ReadFile(Pwd() + \"favicon.ico\") Exception:" + err.Error())
		}
		w.Write(bytes)
		return

	})

	FileServer(r, "/css", http.Dir(Pwd()+"css"))
	FileServer(r, "/js", http.Dir(Pwd()+"js"))

	r.Post("/aes_decrypt", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Println("r.ParseMultipartForm(32 << 20) Exception: " + err.Error())
			w.WriteHeader(400)
			w.Write([]byte("r.ParseMultipartForm(32 << 20) Exception: " + err.Error()))
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Println("r.FormFile(\"file\") Exception: " + err.Error())
			w.WriteHeader(400)
			w.Write([]byte("r.FormFile(\"file\") Exception: " + err.Error()))
			return
		}

		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Println("io.Copy(buf, file) Exception: " + err.Error())
			w.WriteHeader(500)
			w.Write([]byte("io.Copy(buf, file) Exception: " + err.Error()))
			return
		}

		decrypted, err := crypto.AESDecryptBytes(buf.Bytes())
		if err != nil {
			log.Println("crypto.AESDecryptBytes(encryptionKey, buf.Bytes()) Exception: " + err.Error())
			w.WriteHeader(400)
			w.Write([]byte("crypto.AESDecryptBytes(encryptionKey, buf.Bytes()) Exception: " + err.Error()))
			return
		}

		w.WriteHeader(200)
		w.Write(decrypted)
	})

	r.Post("/aes_encrypt", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Println("r.ParseMultipartForm(32 << 20) Exception: " + err.Error())
			w.WriteHeader(400)
			w.Write([]byte("r.ParseMultipartForm(32 << 20) Exception: " + err.Error()))
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Println("r.FormFile(\"file\") Exception: " + err.Error())
			w.WriteHeader(400)
			w.Write([]byte("r.FormFile(\"file\") Exception: " + err.Error()))
			return
		}

		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Println("io.Copy(buf, file) Exception: " + err.Error())
			w.WriteHeader(500)
			w.Write([]byte("io.Copy(buf, file) Exception: " + err.Error()))
			return
		}

		encrypted, err := crypto.AESEncryptBytes(buf.Bytes())
		if err != nil {
			log.Println("crypto.AESEncryptBytes(encryptionKey, buf.Bytes()) Exception: " + err.Error())
			w.WriteHeader(400)
			w.Write([]byte("crypto.AESEncryptBytes(encryptionKey, buf.Bytes()) Exception: " + err.Error()))
			return
		}

		w.WriteHeader(200)
		w.Write(encrypted)
	})

	go func() {
		time.Sleep(5 * time.Millisecond)
		fmt.Println("Running Go AES Program on port: " + PORT)
	}()

	log.Fatal(http.ListenAndServe(":"+PORT, r))

}

func HTTPFile(path string, w http.ResponseWriter) {
	file, _ := os.Open(path)
	defer file.Close()
	fileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	file.Read(fileHeader)
	//Get content type of file
	fileStat, _ := file.Stat()
	fileSize := strconv.Itoa(int(fileStat.Size()))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", fileSize)
	file.Seek(0, 0)
	io.Copy(w, file) //'Copy' the file to the client
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func Pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return dir + "/"
}
