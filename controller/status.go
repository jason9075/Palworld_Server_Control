package controller

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
    "encoding/json"
)

type ServerStatus struct{
    Name string `json:"name"`
    MaxPlayers int `json:"maxplayers"`
    NumPlayers int `json:"numplayers"`
    Players []string `json:"players"`
    Raw ServerRaw `json:"raw"`
}

type ServerRaw struct{
    Started bool `json:"started"`
    Attributes ServerAttributes `json:"attributes"`
}

type ServerAttributes struct{
    Players int `json:"PLAYERS_l"`
    Days int `json:"DAYS_l"`
}

type ServerError struct{
    Error string `json:"error"`
}

type Response struct{
    Running bool `json:"running"`
    Status ServerStatus `json:"status"`
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

    var response Response
    w.Header().Set("Content-Type", "application/json")

    var error ServerError
    err = json.Unmarshal(output, &error)
    if err == nil && error.Error != "" {
        response.Running = false
        json.NewEncoder(w).Encode(response)
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
    response.Running = true
    response.Status = status
    json.NewEncoder(w).Encode(response)

}
