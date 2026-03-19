import Phaser from 'phaser'

const GRID_SIZE = 8
const TILE_W = 128
const TILE_H = 64

const BUILDING_TEXTURES: Record<string, string> = {
  hq: 'building-hq',
  barracks: 'building-barracks',
  factory: 'building-factory',
  power: 'building-power',
  turret: 'building-turret',
  radar: 'building-radar',
}

export class BaseScene extends Phaser.Scene {
  private buildings: { id: string; x: number; y: number; sprite: Phaser.GameObjects.Image }[] = []
  private selectedBuilding: string | null = null
  private coinsText!: Phaser.GameObjects.Text
  private coins = 500

  constructor() {
    super({ key: 'BaseScene' })
  }

  create() {
    const centerX = this.cameras.main.width / 2
    const centerY = 200

    // Draw isometric grid
    for (let row = 0; row < GRID_SIZE; row++) {
      for (let col = 0; col < GRID_SIZE; col++) {
        const { x, y } = this.toIso(row, col, centerX, centerY)
        const tile = this.add.image(x, y, 'tile-grass')
        tile.setInteractive()
        tile.setData('row', row)
        tile.setData('col', col)
        tile.on('pointerdown', () => this.onTileClick(row, col, centerX, centerY))
        tile.on('pointerover', () => tile.setTint(0x44ff44))
        tile.on('pointerout', () => tile.clearTint())
      }
    }

    // Place HQ at center
    this.placeBuilding('hq', 3, 3, centerX, centerY)

    // UI
    this.createUI()
  }

  toIso(row: number, col: number, offsetX: number, offsetY: number) {
    return {
      x: offsetX + (col - row) * (TILE_W / 2),
      y: offsetY + (col + row) * (TILE_H / 2),
    }
  }

  createUI() {
    const w = this.cameras.main.width

    // Top bar
    this.add.rectangle(w / 2, 30, w, 60, 0x111118, 0.9).setDepth(100)
    this.add.text(20, 18, 'COMMAND CENTER', { fontSize: '18px', color: '#8b5cf6', fontStyle: 'bold' }).setDepth(101)
    this.coinsText = this.add.text(w - 200, 18, `Coins: ${this.coins}`, { fontSize: '16px', color: '#f59e0b' }).setDepth(101)

    // Building palette
    const paletteY = this.cameras.main.height - 100
    this.add.rectangle(w / 2, paletteY + 20, w, 120, 0x111118, 0.9).setDepth(100)
    this.add.text(20, paletteY - 20, 'BUILD:', { fontSize: '12px', color: '#6b7280' }).setDepth(101)

    const buildOptions = [
      { id: 'barracks', name: 'Barracks', cost: 100, color: '#3b82f6' },
      { id: 'factory', name: 'Factory', cost: 250, color: '#6366f1' },
      { id: 'power', name: 'Power', cost: 75, color: '#f59e0b' },
      { id: 'turret', name: 'Turret', cost: 150, color: '#ef4444' },
      { id: 'radar', name: 'Radar', cost: 200, color: '#10b981' },
    ]

    buildOptions.forEach((opt, i) => {
      const x = 80 + i * 140
      const btn = this.add.rectangle(x, paletteY + 10, 120, 50, 0x1f1f2e).setDepth(101).setInteractive()
      this.add.text(x - 50, paletteY - 5, opt.name, { fontSize: '13px', color: opt.color }).setDepth(102)
      this.add.text(x - 50, paletteY + 15, `${opt.cost} coins`, { fontSize: '11px', color: '#6b7280' }).setDepth(102)

      btn.on('pointerdown', () => {
        this.selectedBuilding = this.selectedBuilding === opt.id ? null : opt.id
        this.updatePaletteSelection(buildOptions, i)
      })
      btn.on('pointerover', () => btn.setFillStyle(0x2a2a3e))
      btn.on('pointerout', () => {
        btn.setFillStyle(this.selectedBuilding === opt.id ? 0x3a3a5e : 0x1f1f2e)
      })
    })

    // Battle button
    const battleBtn = this.add.rectangle(w - 100, paletteY + 10, 140, 50, 0x8b5cf6).setDepth(101).setInteractive()
    this.add.text(w - 155, paletteY + 2, 'ATTACK!', { fontSize: '16px', color: '#fff', fontStyle: 'bold' }).setDepth(102)
    battleBtn.on('pointerdown', () => this.scene.start('BattleScene'))
    battleBtn.on('pointerover', () => battleBtn.setFillStyle(0x7c3aed))
    battleBtn.on('pointerout', () => battleBtn.setFillStyle(0x8b5cf6))
  }

  updatePaletteSelection(_options: { id: string }[], _selectedIndex: number) {
    // Visual feedback handled by pointerout
  }

  onTileClick(row: number, col: number, centerX: number, centerY: number) {
    if (!this.selectedBuilding) return

    // Check if tile is occupied
    if (this.buildings.some((b) => b.x === row && b.y === col)) return

    const costs: Record<string, number> = { barracks: 100, factory: 250, power: 75, turret: 150, radar: 200 }
    const cost = costs[this.selectedBuilding] ?? 0
    if (this.coins < cost) return

    this.coins -= cost
    this.coinsText.setText(`Coins: ${this.coins}`)
    this.placeBuilding(this.selectedBuilding, row, col, centerX, centerY)
    this.selectedBuilding = null
  }

  placeBuilding(id: string, row: number, col: number, centerX: number, centerY: number) {
    const { x, y } = this.toIso(row, col, centerX, centerY)
    const texture = BUILDING_TEXTURES[id] ?? 'building-hq'
    const sprite = this.add.image(x, y - 30, texture).setDepth(row + col + 10)
    this.buildings.push({ id, x: row, y: col, sprite })

    // Place animation
    sprite.setScale(0)
    this.tweens.add({ targets: sprite, scaleX: 1, scaleY: 1, duration: 300, ease: 'Back.easeOut' })
  }
}
