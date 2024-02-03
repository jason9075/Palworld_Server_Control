package main

import (
    "net/http"
    "fmt"
    "pws/controller"
    "github.com/joho/godotenv"
)


func main() {
    err := godotenv.Load()
        if err != nil {
        fmt.Println("Error loading .env file")
    }

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)
    http.HandleFunc("/startServer", controller.StartServerHandler)
    http.HandleFunc("/stopServer", controller.StopServerHandler)
    http.HandleFunc("/fetchServerConfig", controller.FetchServerConfigHandler)
    http.HandleFunc("/setServerConfig", controller.SetServerConfigHandler)
    http.HandleFunc("/status", controller.StatusHandler)

    fmt.Println("Server is running on port 8080...")
    http.ListenAndServe(":8080", nil)
}

