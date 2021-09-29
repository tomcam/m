package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type app struct {
	errLog  *log.Logger
	infoLog *log.Logger
}

func main() {
  port := flag.String("port", "12345", "Port to run on")
  flag.Parse()

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  app := &app{
    errLog: errLog,
    infoLog: infoLog,
  }

  mux := http.NewServeMux()
  mux.HandleFunc("/", app.root)
  mux.HandleFunc("/new", app.newPage)
  
  //fs := http.FileServer(http.Dir("./ui/html/"))
  //mux.Handle("/static/", http.StripPrefix("/static", fs))

  site := &http.Server{
    Addr: ":" + *port,
    ErrorLog: errLog,
    Handler: mux,
  }

  infoLog.Printf("Running on port %s. Press Ctrl+C to stop", *port)
  err := site.ListenAndServe()
  errLog.Fatal(err)
}




