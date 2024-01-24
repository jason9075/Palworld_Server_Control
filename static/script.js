function startServer() {
  var password = document.getElementById('passwordInput').value
  fetch('/startServer', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ password: password }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error('Password incorrect or server error')
      }
      alert('Server started, please wait 30s for it to be ready')
      return response.text()
    })
    .then((data) => {
      document.getElementById('error').innerText = ''
      updateStatus()
    })
    .catch((error) => {
      document.getElementById('error').innerText = error.message
    })
}

function stopServer() {
  var password = document.getElementById('passwordInput').value
  fetch('/stopServer', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ password: password }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error('Password incorrect or server error')
      }
      return response.text()
    })
    .then((data) => {
      document.getElementById('error').innerText = ''
      alert('Server stopped')
      updateStatus()
    })
    .catch((error) => {
      document.getElementById('error').innerText = error.message
    })
}

function updateStatus() {
  fetch('/status')
    .then((response) => response.json())
    .then((data) => {
      const coverImg = document.querySelector('.cover-img')
      if (data.running) {
        document.getElementById('status').innerText =
          `遊戲伺服器正在運行\n 玩家人數: ${data.numPlayers}/${data.maxPlayers}\n 玩家ID: ${data.players}`
        coverImg.src = 'images/cover.jpg'
      } else {
        document.getElementById('status').innerText = '遊戲伺服器未運行'
        coverImg.src = 'images/cover-sleep.png'
      }
    })
    .catch((error) => console.error('Error:', error))
}

setInterval(updateStatus, 10000)
updateStatus()
