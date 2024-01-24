package controller

import (
    "net/http"
    "golang.org/x/crypto/ssh"
    "encoding/json"
    "fmt"
    "time"
    "os"
    "bytes"
    "os/exec"
)

type PasswdRequest struct {
    Password string `json:"password"`
}

func StartServerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var pass PasswdRequest
    err := json.NewDecoder(r.Body).Decode(&pass)
    if err != nil {
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    if pass.Password != os.Getenv("PASSWORD") {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // send wakelan msg to device
    cmd := exec.Command(os.Getenv("LOCAL_WAKE_SCRIPT"))
    err = cmd.Run()
    if err != nil {
        fmt.Println("Error executing wake script:", err)
    }

    err = executeSSHCommand(os.Getenv("START_COMMAND"), os.Getenv("SERVER_HOST"), os.Getenv("SSH_USER"), os.Getenv("SSH_KEY_PATH"))
    if err != nil {
        fmt.Println("Error executing SSH command:", err)
    } else {
        discordMsg("帕爾，啟動！")
        fmt.Println("Command executed successfully")
    }

}

func StopServerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var pass PasswdRequest
    err := json.NewDecoder(r.Body).Decode(&pass)
    if err != nil {
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    if pass.Password != os.Getenv("PASSWORD") {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    err = executeSSHCommand(os.Getenv("STOP_COMMAND"), os.Getenv("SERVER_HOST"), os.Getenv("SSH_USER"), os.Getenv("SSH_KEY_PATH"))
    if err != nil {
        fmt.Println("Error executing SSH command:", err)
    } else {
        fmt.Println("Command executed successfully")
    }

}

func executeSSHCommand(command, hostname, username, keyPath string) error {
    key, err := os.ReadFile(keyPath)
    if err != nil {
        return fmt.Errorf("unable to read private key: %v", err)
    }

    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        return fmt.Errorf("unable to parse private key: %v", err)
    }

    config := &ssh.ClientConfig{
        User: username,
        Auth: []ssh.AuthMethod{
            ssh.PublicKeys(signer),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 請注意這是不安全的，僅用於測試
        Timeout:         10 * time.Second,
    }

    client, err := ssh.Dial("tcp", hostname+":22", config)
    if err != nil {
        return fmt.Errorf("unable to connect: %v", err)
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        return fmt.Errorf("unable to create session: %v", err)
    }
    defer session.Close()

    err = session.Run(command)
    if err != nil {
        return fmt.Errorf("unable to run command: %v", err)
    }

    return nil
}

func discordMsg(msg string) {
    webhookUrl := os.Getenv("DISCORD_WEBHOOK")
    if webhookUrl == "" {
        return
    }

    jsonValue := []byte(fmt.Sprintf(`{"content":"%s"}`, msg))
    _, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
        fmt.Println("Error sending discord message:", err)
    }
}
