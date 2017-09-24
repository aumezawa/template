package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "html/template"
  "net/http"
)

const MAX_FROM_BUFFER_SIZE = 2 * 1024 * 1024
const MAX_FILE_BUFFER_SIZE = 1 * 1024 * 1024

var (
  optPort = flag.Int("p", 80, "HTTP Port Number")
  optHtmlFile = flag.String("f", "uploadFile.html", "Index HTML file")
)

var (
  port int
  htmlFile string
)

type Content struct {
  Title string
  URI string
  ContentTitle string
  ContentMain string
}

func Panic(err error) {
  if err != nil {
    panic(err)
  }
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles(htmlFile)
  Panic(err)

  switch r.Method {
  case "GET":
    content := Content{"Upload File", "", "", ""}
    err = tmpl.Execute(w, content)
    Panic(err)

  case "POST":
    r.ParseMultipartForm(MAX_FROM_BUFFER_SIZE)
    file, handler, err := r.FormFile("uploadfile")

    buf := make([]byte, MAX_FILE_BUFFER_SIZE)
    n, err := file.Read(buf)
    Panic(err)
    defer file.Close()

    content := Content{"Upload File", "", handler.Filename, fmt.Sprintf("%s", string(buf[:n]))}
    err = tmpl.Execute(w, content)
    Panic(err)

  default:
  }

}

func main() {
  flag.Parse()

  port = *optPort
  if port < 0 || port > 65535 {
    log.Fatal("Invalid port number")
  }

  htmlFile = *optHtmlFile
  if _, err := os.Stat(htmlFile); err != nil {
    log.Fatal("Not exist the HTML file")
  }

  log.Printf("Starting web service on port %d", port)
  http.HandleFunc("/", uploadFileHandler)
  http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
