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
    Payload     string `json:"payload"`
}
func FetchServerConfigHandler(w http.ResponseWriter, r *http.Request) {
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

    config, err := executeSSHCommand(os.Getenv("FETCH_CONFIG_COMMAND"), os.Getenv("SERVER_HOST"), os.Getenv("SSH_USER"), os.Getenv("SSH_KEY_PATH"))
    if err != nil {
        fmt.Println("Error executing SSH command:", err)
        return
    }

    w.Write([]byte(config))
}

func SetServerConfigHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var data PasswdRequest
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    if data.Password != os.Getenv("PASSWORD") {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    _, err = executeSSHCommand(os.Getenv("SET_CONFIG_COMMAND")+" '"+data.Payload+"'", os.Getenv("SERVER_HOST") , os.Getenv("SSH_USER"), os.Getenv("SSH_KEY_PATH"))
    if err != nil {
        fmt.Println("Error executing SSH command:", err)
        return
    }

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

    _, err = executeSSHCommand(os.Getenv("START_COMMAND"), os.Getenv("SERVER_HOST"), os.Getenv("SSH_USER"), os.Getenv("SSH_KEY_PATH"))
    if err != nil {
        fmt.Println("Error executing SSH command:", err)
    } else {
        discordMsg(randomMsg())
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

    _, err = executeSSHCommand(os.Getenv("STOP_COMMAND"), os.Getenv("SERVER_HOST"), os.Getenv("SSH_USER"), os.Getenv("SSH_KEY_PATH"))
    if err != nil {
        fmt.Println("Error executing SSH command:", err)
    } else {
        fmt.Println("Command executed successfully")
    }

}

func executeSSHCommand(command, hostname, username, keyPath string) (string, error) {
    key, err := os.ReadFile(keyPath)
    if err != nil {
        return "", fmt.Errorf("unable to read private key: %v", err)
    }

    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        return "", fmt.Errorf("unable to parse private key: %v", err)
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
        return "", fmt.Errorf("unable to connect: %v", err)
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        return "", fmt.Errorf("unable to create session: %v", err)
    }
    defer session.Close()

    var stdoutBuf bytes.Buffer
    session.Stdout = &stdoutBuf

    err = session.Run(command)
    if err != nil {
        return "", fmt.Errorf("unable to run command: %v", err)
    }

    return stdoutBuf.String(), nil
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

func randomMsg() string {
    msgs := []string{
        "帕魯，啟動！",
        "早上好！帕魯，現在是工作時間",
        "我來到一個島 它叫帕魯奴隸島",
        "你不幹有的是帕魯幹",
        "帕魯，心中最軟的一愧",
        "我帕魯諾•喬巴納有個夢想，那就是...",
        "今天玩帕魯，明天當帕魯",
        "Dio：帕魯，瓦魯多",
    }

    return msgs[time.Now().UnixNano() % int64(len(msgs))]
}
