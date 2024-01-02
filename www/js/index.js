document.addEventListener('DOMContentLoaded', () => {
  bindElements(document)
  if (isLoggedIn()) {
    showGameList()
  } else {
    hideAllMainsExcept('login')
  }
})

function isLoggedIn() {
  return localStorage.getItem('username') && localStorage.getItem('userid')
}

function bindElements(root) {
  gameList = root.getElementById('games')
  root.getElementById('form-login').addEventListener('submit', xhrFormSubmit)
  root.getElementById('form-login').addEventListener('formdata', loginFormData)

  root.getElementById('form-game-create').addEventListener('submit', xhrFormSubmit)
  root.getElementById('form-game-create').addEventListener('formdata', createGameData)

  root.getElementById('refresh-button').addEventListener('click', () => RefreshGameList(gameList))
  root.getElementById('create-game-button').addEventListener('click', showCreateGameScreen)
  root.getElementById('exit-new-game').addEventListener('click', showGameList)
}

function xhrFormSubmit(evt) {
  evt.preventDefault()
  new FormData(evt.target)
}

function loginFormData(evt) {
  username = evt.formData.get('userName')
  fetch('/users', {
    method: 'POST',
    body: JSON.stringify({ name: username }),
  }).then(res => res.json())
    .then(data => {
      localStorage.setItem('username', data.name)
      localStorage.setItem('userid', data.id)
      showGameList()
    })
    .catch(err => { console.warn(err) })
}

function createGameData(evt) {
  gameData = {
    name: evt.formData.get('gameName')
  }
  fetch('/games', {
    method: 'POST',
    body: JSON.stringify(gameData),
    headers: {
      'X-Player-Id': localStorage.getItem('userid'),
    }
  })
    .then(res => res.json())
    .then(data => { console.log(data) })
    .catch(err => { console.warn(err) })
}

function CreateGameListEntry(game) {
  const li = document.createElement('li')

  const owner = document.createElement('strong')
  owner.innerText = game.owner

  const button = document.createElement('input')
  button.setAttribute('type', 'button')
  button.setAttribute('value', 'Show')
  button.addEventListener('click', () => fetch(`/games/${game.id}`)
    .then(res => {
      if (res.ok)
        return res.json()
      throw new Error('Error: ' + res.status)
    })
    .then(showGame)
    .catch(err => { console.warn(err) })
  )

  const name = document.createElement('span')
  name.setAttribute('class', 'game-name')
  name.innerText = game.name

  li.appendChild(owner)
  li.appendChild(name)
  li.appendChild(button)
  return li
}

function RefreshGameList(gl) {
  fetch('/games')
    .then(resp => {
      if (resp.ok)
        return resp.json()
      throw new Error('Error: ' + resp.status)
    })
    .then(data => {
      while (gl.lastChild)
        gl.removeChild(gl.lastChild)
      data
        .map(CreateGameListEntry)
        .forEach(item => gl.append(item))
    })
    .catch(err => { console.warn(err) })
}

function hideAllMainsExcept(id) {
  const mains = document.getElementsByTagName('main')
  for (let i = 0; i < mains.length; i++)
    mains[i].hidden = !(mains[i].id === id)
}

function showCreateGameScreen() {
  hideAllMainsExcept('new-game')
}

function showGameList() {
  RefreshGameList(gameList)
  hideAllMainsExcept('game-list')
}

function createGame() {
  newGameData = {
    name: document.getElementById('input-text-game-name').value
  }
  fetch('/games', {
    method: 'POST',
    headers: {
      'X-Player-Id': localStorage.getItem('userid'),
    },
    body: JSON.stringify(newGameData),
  }).then(res => res.json())
    .then(data => showGame(data))
    .catch(err => { console.warn(err) })
}

function showGame(data) {
  updateGameDisplay(data)
  hideAllMainsExcept('game')
}

function updateGameDisplay(data) {
  document.getElementById('game-name').innerText = data.name
  const playerList = document.getElementById('game-players')
  data.playerNames.forEach(p => {
    const e = document.createElement('li')
    e.textContent = p
    playerList.append(e)
  })
}

function createUser() {
  userData = {
    name: document.getElementById('input-text-user-name').value
  }
  fetch('/users', {
    method: 'POST',
    body: JSON.stringify(userData),
  }).then(res => res.json())
    .then((data) => {
      localStorage.setItem('username', data.name)
      localStorage.setItem('userid', data.id)
      showGameList()
    })
    .catch(err => { console.warn(err) })
}

function previewGame(id) {
  hideAllMainsExcept('game-preview')
  fetch(`/games/${id}`, {
    headers: {
      'Accept': 'application/json',
    },
  })
    .then(res => res.json())
    .then(json => {
      document.getElementById('game-name').innerText = json.name
      document.getElementById('game-players').innerText = json.players
    })
    .catch(err => console.warn(err))
}

