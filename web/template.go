package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "html/template"
  "net/http"
)

var (
  optPort = flag.Int("p", 80, "HTTP Port Number")
  optHtmlFile = flag.String("f", "template.html", "Index HTML file")
)

var (
  port int
  htmlFile string
)

type Content struct {
  Title string
  Body string
}

func Panic(err error) {
  if err != nil {
    panic(err)
  }
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles(htmlFile)
  Panic(err)

  switch r.Method {
  case "GET":
    content := Content{"Template Title", "Template Body"}
    err = tmpl.Execute(w, content)
    Panic(err)

  case "POST":

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
  http.HandleFunc("/", templateHandler)
  http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
