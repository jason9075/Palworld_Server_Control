package controller

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
    serverAddress := os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
    if checkServerStatus(serverAddress) {
        fmt.Fprint(w, "--")
        // fmt.Fprint(w, "遊戲伺服器正在運行")
    } else {
        fmt.Fprint(w, "-")
        // fmt.Fprint(w, "遊戲伺服器未運行")
    }
}

func checkServerStatus(address string) bool {
    timeout := 10 * time.Second
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        fmt.Println(err.Error())
        return strings.Contains(err.Error(), "connection refused")
    }
    defer conn.Close()
    fmt.Println("Connection successful")

    return false
}

