# Palworld Game Server Control Panel

## Introduction

This project is a web-based control panel for starting and stopping a game server. It is designed to run on a Raspberry Pi, acting as a router, with the game server running on a separate Ubuntu machine. The control panel is built using Go and provides a simple web interface for server control.

## Features

- Start and stop the game server remotely.
- Check the current status of the game server.
- Secure server control with password protection.

## Prerequisites

- A Raspberry Pi set up as a router.
- A game server running on Ubuntu.
- Go programming language installed on the Raspberry Pi.
- SSH access from the Raspberry Pi to the Ubuntu server.

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/game-server-control.git
   cd game-server-control
   ```

2. **Set Up the .env File**

   Create a `.env` file in the root directory of the project with the following content:

   ```
   SSH_KEY_PATH=/path/to/your/ssh_key
   SSH_USER=your_ssh_username
   SERVER_HOST=your_server_ip
   SERVER_PORT=your_server_port
   START_COMMAND=path/to/start_script.sh
   STOP_COMMAND=path/to/stop_script.sh
   PASSWORD=your_password
   ```

   Replace the values with your specific server and SSH details.

3. **Install Dependencies**

   ```bash
   go get .
   ```

   ```bash
   npm install gamedig -g
   ```

## Running the Server

To start the web server, run:

```bash
go run main.go
```

The web interface will be accessible at `http://localhost:8080` or the respective IP and port of your Raspberry Pi.

## Using the Control Panel

- Open a web browser and navigate to `http://raspberry_pi_ip:8080`.
- Enter the designated password in the password field.
- Use the "Start Server" and "Stop Server" buttons to control the game server.

## Security

Ensure that the Raspberry Pi and the game server have firewalls enabled and are configured to allow only necessary traffic. The use of SSH keys and secure password practices is recommended.

---

## License

[MIT License](LICENSE)
