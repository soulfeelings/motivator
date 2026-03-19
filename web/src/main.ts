import Phaser from 'phaser'
import { BootScene } from './scenes/BootScene'
import { BaseScene } from './scenes/BaseScene'
import { BattleScene } from './scenes/BattleScene'

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  parent: 'game-container',
  width: window.innerWidth,
  height: window.innerHeight,
  backgroundColor: '#0a0a0f',
  scene: [BootScene, BaseScene, BattleScene],
  scale: {
    mode: Phaser.Scale.RESIZE,
    autoCenter: Phaser.Scale.CENTER_BOTH,
  },
  pixelArt: false,
  antialias: true,
}

new Phaser.Game(config)
