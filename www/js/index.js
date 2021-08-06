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
	root.getElementById('refresh-button').addEventListener('click', (evt) => RefreshGameList(gameList))
	root.getElementById('create-game-button').addEventListener('click', showCreateGameScreen)
	root.getElementById('exit-new-game').addEventListener('click', showGameList)
	root.getElementById('input-button-create-game').addEventListener('click', createGame)
	root.getElementById('input-button-create-user').addEventListener('click', createUser)
}

function CreateGameListEntry(game) {
	const li = document.createElement('li')
	li.textContent = game.name + ' ' + game.players
	return li
}

function RefreshGameList(gl) {
	fetch('http://localhost:8080/games')
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
}

function hideAllMainsExcept(id) {
	const mains = document.getElementsByTagName('main')
	for (let i=0; i < mains.length; i++)
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
		name: document.getElementById("input-text-game-name").value
	}
	fetch('/games', {
		method: 'POST',
		body: JSON.stringify(newGameData),
	}).then(res => res.json())
		.then(data => showGame(data))
}

function showGame(data) {
	updateGameDisplay(data)
	hideAllMainsExcept('game')
}

function updateGameDisplay(data) {
	document.getElementById('game-name').innerText = data.name
	const playerList = document.getElementById('game-players')
	data.players.forEach(p => {
		const e = document.createElement('li')
		e.textContent = p
		playerList.append(e)
	})
}

function createUser() {
	userData = {
		name: document.getElementById("input-text-user-name").value
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
}
