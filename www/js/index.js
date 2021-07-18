document.addEventListener('DOMContentLoaded', () => {
	const gameList = document.getElementById('games')
	document.getElementById('refresh-button').addEventListener('click', (evt) => RefreshGameList(gameList))
	document.getElementById('create-game-button').addEventListener('click', showCreateGameScreen)
	document.getElementById('exit-new-game').addEventListener('click', showGameList)
	RefreshGameList(gameList)
})

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
	hideAllMainsExcept('game-list')
}
