<template>
  <div>
    <h1>Game!</h1>
    <div id="game">
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
    components: {
    	ConsoleComponent
    }
  }
</script>