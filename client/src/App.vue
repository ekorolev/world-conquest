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
    <div id="game"></div>
    <div>
    	<ul>
    		<li v-for="message in messages">{{ message }}</li>
    	</ul>
    </div>
    <console-component></console-component>
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

	function preload () {
	}

	function create () {
	}

	function update () {
	}

	import ConsoleComponent from './ConsoleComponent.vue'

  export default {
    data () {
    	return {
    		name: null,
    		auth: false,
    		socket: null,
    		messages: []
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
    			self.messages.push(evt.data)
    		}
    	}
    },
    components: {
    	ConsoleComponent
    }
  }
</script>