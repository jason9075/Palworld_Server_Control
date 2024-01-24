package controller

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
    "encoding/json"
)

type ServerStatus struct{
    Name string `json:"name"`
    MaxPlayers int `json:"maxplayers"`
    NumPlayers int `json:"numplayers"`
    Players []string `json:"players"`
}

type ServerError struct{
    Error string `json:"error"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
    serverAddress := os.Getenv("PUBLIC_SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
    cmd := exec.Command("gamedig", "--type", "palworld", serverAddress)
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println("Error executing gamedig command:", err)
        fmt.Println(string(output))
        http.Error(w, "Error executing gamedig command", http.StatusInternalServerError)
        return
    }
    fmt.Println(string(output))

    var error ServerError
    err = json.Unmarshal(output, &error)
    if err == nil && error.Error != "" {
        fmt.Fprint(w, "遊戲伺服器未運行")
        return
    }


    var status ServerStatus
    err = json.Unmarshal(output, &status)
    if err != nil {
        fmt.Println("Error unmarshalling gamedig output:", err)
        fmt.Println(string(output))
        http.Error(w, "Error unmarshalling gamedig output", http.StatusInternalServerError)
        return
    }
    fmt.Println(status)
    fmt.Fprintf(w, "遊戲伺服器正在運行\n 玩家人數: %d/%d\n 玩家ID: %v", status.NumPlayers, status.MaxPlayers, status.Players)
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

