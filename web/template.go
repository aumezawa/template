package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "path"
  "text/template"
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

func templateHandler(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles(htmlFile)
  if err != nil {
    log.Printf(err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  switch r.Method {
  case "GET":
    content := Content{"Template Title", "Template Body"}
    err = tmpl.Execute(w, content)
    if err != nil {
      log.Printf(err.Error())
      w.WriteHeader(http.StatusInternalServerError)
    }
    return

  case "POST":
    return

  default:
    return
  }

}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case "GET":
    var content string
    file, err := ioutil.ReadFile(r.URL.Path[1:])
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    switch path.Ext(r.URL.Path[1:])[1:] {
    case "js":
      content = "text/javascript"
    case "css":
      content = "text/css"
    default:
      content = "text/plane"
    }
    w.Header().Set("Content-Type", content)
    w.WriteHeader(http.StatusOK)
    w.Write(file)

  default:
    w.WriteHeader(http.StatusBadRequest)
    return
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
  http.HandleFunc("/lib/", getFileHandler)
  http.HandleFunc("/", templateHandler)
  http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
