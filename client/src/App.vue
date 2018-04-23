<template>
  <div>

    <div class="container">
      <div class="row">
        <div class="col" style="width: 250px">
          <div v-if="auth && player">
            Hello, <b>{{ player.name }}</b>! Your team is <img v-bind:src="`/assets/type-${player.team}.png`" class="flag">
            <br/>You may logout <button v-on:click="logoutAuth()">logout</button>
            <br/>
            You may to change your team:
            <ul class="list-unstyled teams">
              <li v-for="(team, key) in teams" v-bind:class="{ 'chosen': (key == player.team) }" v-on:click="changeTeam(key, team)">
                <img v-bind:src="team.image" class="flag"> {{ team.name }}
              </li>
            </ul>
          </div>
          <div v-else>
            Hello! Enter your name for join the Game!<br>
            You will may to choose favourite country and you will start Great Battle!
          </div>
        </div>
        <div class="col">
          <div v-if="!auth || !player">
            <input type="text" v-model="name" placeholder="Enter your name">
            <input type="checkbox" v-model="remember"> Запомнить меня
            <button v-on:click="authName()">Это я</button>
          </div>
          <div v-else>
            <span v-if="allow==0">Можно захватывать!</span>
            <span v-else>Осталось {{ allow }} секунд</span>
          </div>
          <div id="game">
            <div class="game-row" v-for="(row, keyRow) in map" :key="`row${keyRow}`">
              <div class="game-cell" v-for="(cell, keyCell) in row" :key="`row${keyRow}cell${keyCell}`" v-on:click="setCell(keyRow, keyCell)">
                <img v-bind:src="`/assets/type-${cell.value}.png`">
              </div>
            </div>
          </div>
        </div>
        <div class="col" style="width: 220px">
          Players:
          <li v-for="pl in players">
            {{ pl.name }} <img v-bind:src="teams[pl.team].image" class="flag">
          </li>
        </div>
        <div class="col" style="margin-left: 20px; width: 200px">
          Statistic:
          <li v-for="team in sortTeams">
            <img v-bind:src="team.image" class="flag"> {{team.prop}}%
          </li>
        </div>
        <div class="col" style="width: 400px;">
          <h5>Log:</h5>
          <li v-for="log in sortLogs">
            {{ log.message }}
          </li>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  let socket = new WebSocket('ws://localhost:9001/ws')

  export default {
    data () {
    	return {
    		name: null,
    		auth: false,
    		socket: null,
        gameState: null,
        map: [],
        players: [],
        player: null,
        remember: false,
        teams: {
          0: { name: "Neutral", image: "/assets/type-0.png", prop: 0 },
          1: { name: "Russia", image: "/assets/type-1.png", prop: 0 },
          2: { name: "China", image: "/assets/type-2.png", prop: 0 },
          3: { name: "USA", image: "/assets/type-3.png", prop: 0 },
          4: { name: "UK", image: "/assets/type-4.png", prop: 0 },
          5: { name: "Mongolia", image: "/assets/type-5.png", prop: 0 },
        },
        allow: 0,
        interval: null,
        logs: []
    	}
    },
    computed: {
      sortTeams: function () {
        return Object.values(this.teams).sort( function (a,b) {
          if (a.prop<b.prop) { 
            return 1 
          } else {
            return -1
          }
        })
      },
      sortLogs: function () {
         return this.logs.sort( function (a,b) {
          if (a.created<b.created) { 
            return 1 
          } else {
            return -1
          }
        })       
      }
    },
    methods: {
      addlog(msg) {
        this.logs.push({ created: Date.now(), message: msg })
        if (this.logs.length > 5) {
          this.logs.splice(4, 1)
        }
      },
    	authName () {
        if (this.remember)
          window.localStorage.setItem('player', this.name)
    		this.socket.send(`AuthName:${this.name}`)
    	},
      logoutAuth() {
        window.localStorage.setItem('player', '')
        document.location.reload()
      },
      setMap (map) {
        console.log(map)
        map.forEach( (row,keyRow) => {
          this.map.push([])
          row.forEach( (cell,keyCell) => {
            this.map[keyRow].push({
              x: keyRow,
              y: keyCell,
              value: cell
            })
          })
        })
      },
      updateCell(x, y, t) {
        this.map[x][y].value = t
      },
      newPlayer(data) {
        console.log(data)
        this.players.push(data)
        this.addlog(`New player: ${data.name}`)
      },
      setPlayers(data) {
        console.log(data)
        this.players = Object.values(data)
      },
      deletePlayer(id) {
        let sindex = null
        this.players.forEach( (item, key) =>{ if (id == item.id) sindex = key })
        this.addlog(`Player ${this.players[sindex].name} left`)
        if (sindex) this.players.splice(sindex, 1)
      },
      setCell(x, y) {
        if (this.allow > 0) {
          alert(`Жди еще ${this.allow} секунд`)
          return
        }
        //this.map[x][y].value = this.player.team
        this.socket.send(`MarkCell:${x}:${y}`)
        this.allow = 5
        let self = this
        clearInterval(self.interval)
        self.interval = setInterval( () => {
          if ( self.allow <= 0 ) {
            self.allow = 0
            clearInterval(self.interval)
          }
          self.allow -= 1
          if (self.allow <= 0) self.allow = 0
        }, 1000)
      },
      setAboutInfo(info) {
        this.auth = true
        this.player = info
      },
      changeTeam(key, team) {
        this.player.team = key
        this.socket.send(`ChangeTeam:${key}`)
      },
      setStats(data) {
        let self = this
        data.forEach( (item, key) => {
          self.teams[key].prop = item
        })
      },
      updatePlayer(id, data) {
        let updated = false
        this.players.forEach( item => {
          console.log( item.id ,id)
          if (item.id == id) {
            item.name = data.name
            item.id = data.id
            item.team = data.team
            updated = true
          }
        })
        if (!updated) {
          this.players.push(data)
          this.addlog(`New player: ${data.name}`)
        }
      }
    },
    created () {
    	let self = this
    	self.socket = socket

      let name = window.localStorage.getItem('player')
      if (name) {
        this.name = name
        setTimeout( () => this.authName(), 200 )
      }

      this.socket.onmessage = function (evt) {
        console.log(evt.data)
    		if (evt.data === 'Auth:OK') {
          self.auth = true
					return 			
    		}

    		if ( self.auth ) {
          if (evt.data.indexOf('Map:')+1) {
            self.setMap(JSON.parse(evt.data.substr(4,1000000)))
          } else if ( evt.data.indexOf('NewPlayer:')+1 ) {
            self.newPlayer(evt.data)
          } else if ( evt.data.indexOf('Players:')+1) {
            self.setPlayers(JSON.parse(evt.data.substr(8, 100000000)))
          } else if ( evt.data.indexOf('AboutInfo:')+1) {
            self.setAboutInfo(JSON.parse(evt.data.substr(10, 100000000)))
          } else if ( evt.data.indexOf('UpdatePlayer:')+1) {
            let parts = evt.data.split(':')
            let pID = parts[1]
            parts = parts.slice(2)
            let data = parts.join(':')
            self.updatePlayer(pID, JSON.parse(data))
          } else if ( evt.data.indexOf('UpdateCell:')+1) {
            let parts = evt.data.split(':')
            let x = parseInt(parts[1]), y = parseInt(parts[2]), t = parseInt(parts[3])
            self.updateCell(x, y, t)
          } else if (evt.data.indexOf('Error:')+1) {
            alert(evt.data.substr(6, 100000))
          } else if (evt.data.indexOf('Logout:')+1) {
            self.deletePlayer(evt.data.substr(7, 100000))
          } else if (evt.data.indexOf('Stats:')+1) {
            self.setStats(JSON.parse(evt.data.substr(6, 100000000)))
          }
    		}
    	}
    }
  }
</script>

<style>
  .col {
    margin: 10px;
    padding: 10px;
    border: 1px solid #aaa;
    background: #eee;
  }
  .row .col {display: inline-block; vertical-align: top;}
  #game {
    border: 1px solid black;
    width: 640px;
    height: 200px;
  }
  .game-row {
    display: inline-block;
  }
  .game-row .game-cell {
    width: 32px;
    height: 20px;
    background: #eeff99;
  }
  .game-cell img {
    width: 32px;
    height: 20px;
  }
  .game-cell .neutral {
    width: 32px;
    height: 20px;
    background: #fff;
  }
  img.flag {
    width: 32px;
    height: 20px;
  }
  li.chosen {
    border: 1px solid silver;
    background: #8e8;
  }
  ul.teams {
    list-style: none;
  }
  ul.teams li {
    border: 1px solid black;
    padding: 3px;
    width: 200px;
  }
</style>