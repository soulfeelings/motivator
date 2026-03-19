import Phaser from 'phaser'

export class BootScene extends Phaser.Scene {
  constructor() {
    super({ key: 'BootScene' })
  }

  create() {
    // Generate isometric tile textures procedurally
    this.generateTileTexture('tile-grass', 0x1a2e1a, 0x2d4a2d)
    this.generateTileTexture('tile-concrete', 0x2a2a35, 0x3a3a4a)
    this.generateBuildingTexture('building-hq', 0x8b5cf6, 60, 80)
    this.generateBuildingTexture('building-barracks', 0x3b82f6, 50, 60)
    this.generateBuildingTexture('building-factory', 0x6366f1, 55, 70)
    this.generateBuildingTexture('building-power', 0xf59e0b, 40, 50)
    this.generateBuildingTexture('building-turret', 0xef4444, 35, 55)
    this.generateBuildingTexture('building-radar', 0x10b981, 45, 65)
    this.generateUnitTexture('unit-attacker', 0x3b82f6)
    this.generateUnitTexture('unit-defender', 0xef4444)

    this.scene.start('BaseScene')
  }

  generateTileTexture(key: string, color1: number, color2: number) {
    const g = this.add.graphics()
    g.fillStyle(color1, 1)
    g.beginPath()
    g.moveTo(64, 0)
    g.lineTo(128, 32)
    g.lineTo(64, 64)
    g.lineTo(0, 32)
    g.closePath()
    g.fillPath()

    g.lineStyle(1, color2, 0.5)
    g.beginPath()
    g.moveTo(64, 0)
    g.lineTo(128, 32)
    g.lineTo(64, 64)
    g.lineTo(0, 32)
    g.closePath()
    g.strokePath()

    g.generateTexture(key, 128, 64)
    g.destroy()
  }

  generateBuildingTexture(key: string, color: number, w: number, h: number) {
    const g = this.add.graphics()

    // Top face (isometric)
    g.fillStyle(color, 0.8)
    g.beginPath()
    g.moveTo(w, 0)
    g.lineTo(w * 2, w / 2)
    g.lineTo(w, w)
    g.lineTo(0, w / 2)
    g.closePath()
    g.fillPath()

    // Left face
    g.fillStyle(color, 0.5)
    g.fillRect(0, w / 2, w, h)
    g.beginPath()
    g.moveTo(0, w / 2)
    g.lineTo(w, w)
    g.lineTo(w, w + h)
    g.lineTo(0, w / 2 + h)
    g.closePath()
    g.fillPath()

    // Right face
    g.fillStyle(color, 0.3)
    g.beginPath()
    g.moveTo(w, w)
    g.lineTo(w * 2, w / 2)
    g.lineTo(w * 2, w / 2 + h)
    g.lineTo(w, w + h)
    g.closePath()
    g.fillPath()

    g.generateTexture(key, w * 2, w + h)
    g.destroy()
  }

  generateUnitTexture(key: string, color: number) {
    const g = this.add.graphics()
    g.fillStyle(color, 1)
    g.fillCircle(8, 8, 8)
    g.fillStyle(0xffffff, 0.4)
    g.fillCircle(6, 6, 3)
    g.generateTexture(key, 16, 16)
    g.destroy()
  }
}
