<template>
  <div>
    <h1>Game!</h1>
    <div v-if="!auth">
    	<input type="text" v-model="name" placeholder="Enter your name">
    	<button v-on:click="authName()">Это я</button>
    </div>
    <div v-else>
    	Залогинен под именем {{ name }}
    </div>
    <div>{{ gameState }}</div>
    <div id="game"></div>
    <div>
    	<ul>
    		<li v-for="message in messages">{{ message }}</li>
    	</ul>
    </div>
  </div>
</template>

<script>
	require('./phaser.min.js')

	var config = {
	    type: Phaser.AUTO,
	    width: 800,
	    height: 300,
	    parent: "game",
	    scene: {
	        preload: preload,
	        create: create,
	        update: update
	    }
	}

	var game = new Phaser.Game(config);
  var Soldier = null

	function preload () {
    this.load.image('soldier', '/assets/soldier_main.gif')
	}

	function create () {
    Soldier = this.add.image(100, 100, 'soldier')
    Soldier.setDisplaySize( 64, 64 )
    Soldier.setRotation( 30 )
	}

	function update () {
    if ( Phaser.Input.Keyboard.JustDown(this.cursors.left) ) {
      Soldier.setX(Soldier.X - 20)
    }
	}

  export default {
    data () {
    	return {
    		name: null,
    		auth: false,
    		socket: null,
    		messages: [],
        gameState: null
    	}
    },
    methods: {
    	authName () {
    		this.socket.send(`AuthName:${this.name}`)
    	}
    },
    created () {
    	let self = this
    	this.socket = new WebSocket('ws://localhost:9001/ws')
    	this.socket.onmessage = function (evt) {
    		if (evt.data === 'Auth:OK') {
					self.auth = true
					return 			
    		}

    		if ( self.auth ) {

          if ( evt.data.indexOf('State:')+1 ) {
            let array = evt.data.split('')
            array.splice(0, 6)
            let stateJSON = array.join('')
            self.gameState = JSON.parse(stateJSON)
            return
          }

          if ( evt.data.indexOf('NewPlayer:')+1 ) {
            self.messages.push(evt.data)
            return
          }
    		}
    	}
    }
  }
</script>