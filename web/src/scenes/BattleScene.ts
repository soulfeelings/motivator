import Phaser from 'phaser'

export class BattleScene extends Phaser.Scene {
  private attackers: Phaser.GameObjects.Image[] = []
  private defenders: Phaser.GameObjects.Image[] = []
  private tick = 0
  private battleOver = false

  constructor() {
    super({ key: 'BattleScene' })
  }

  create() {
    const w = this.cameras.main.width
    const h = this.cameras.main.height

    // Background
    this.add.rectangle(w / 2, h / 2, w, h, 0x0a0a0f)

    // Battlefield grid
    const centerX = w / 2
    const centerY = h / 2 - 50
    for (let row = 0; row < 10; row++) {
      for (let col = 0; col < 12; col++) {
        const x = centerX + (col - row) * 32 - 100
        const y = centerY + (col + row) * 16 - 80
        this.add.image(x, y, 'tile-concrete').setScale(0.5).setAlpha(0.3)
      }
    }

    // Title
    this.add.text(w / 2 - 80, 20, 'BATTLE IN PROGRESS', { fontSize: '20px', color: '#ef4444', fontStyle: 'bold' })

    // Spawn units
    for (let i = 0; i < 8; i++) {
      const ax = 100 + Math.random() * 150
      const ay = h / 2 - 100 + Math.random() * 200
      const attacker = this.add.image(ax, ay, 'unit-attacker').setScale(1.5).setDepth(50)
      this.attackers.push(attacker)
    }

    for (let i = 0; i < 6; i++) {
      const dx = w - 100 - Math.random() * 150
      const dy = h / 2 - 100 + Math.random() * 200
      const defender = this.add.image(dx, dy, 'unit-defender').setScale(1.5).setDepth(50)
      this.defenders.push(defender)
    }

    // Back button
    const backBtn = this.add.rectangle(80, h - 40, 120, 40, 0x1f1f2e).setDepth(100).setInteractive()
    this.add.text(35, h - 50, 'Back to Base', { fontSize: '13px', color: '#9ca3af' }).setDepth(101)
    backBtn.on('pointerdown', () => this.scene.start('BaseScene'))

    // Start battle simulation
    this.time.addEvent({ delay: 800, callback: () => this.simulateTick(), loop: true })
  }

  simulateTick() {
    if (this.battleOver) return
    this.tick++

    const alive = (units: Phaser.GameObjects.Image[]) => units.filter((u) => u.active)
    const aliveAttackers = alive(this.attackers)
    const aliveDefenders = alive(this.defenders)

    if (aliveAttackers.length === 0 || aliveDefenders.length === 0 || this.tick > 20) {
      this.battleOver = true
      const won = aliveAttackers.length > aliveDefenders.length
      const w = this.cameras.main.width
      const h = this.cameras.main.height

      this.add.rectangle(w / 2, h / 2, 400, 120, 0x111118, 0.95).setDepth(200)
      this.add.rectangle(w / 2, h / 2, 400, 120, won ? 0x10b981 : 0xef4444, 0.3).setDepth(200)
      this.add.text(w / 2 - 60, h / 2 - 30, won ? 'VICTORY!' : 'DEFEAT', {
        fontSize: '32px', color: won ? '#10b981' : '#ef4444', fontStyle: 'bold',
      }).setDepth(201)
      this.add.text(w / 2 - 50, h / 2 + 15, won ? '+50 coins, +100 XP' : 'Better luck next time', {
        fontSize: '14px', color: '#9ca3af',
      }).setDepth(201)
      return
    }

    // Each attacker moves toward nearest defender and attacks
    for (const attacker of aliveAttackers) {
      const target = this.findNearest(attacker, aliveDefenders)
      if (!target) continue

      const dist = Phaser.Math.Distance.Between(attacker.x, attacker.y, target.x, target.y)
      if (dist > 30) {
        // Move toward target
        const angle = Phaser.Math.Angle.Between(attacker.x, attacker.y, target.x, target.y)
        this.tweens.add({
          targets: attacker, x: attacker.x + Math.cos(angle) * 20, y: attacker.y + Math.sin(angle) * 20,
          duration: 600, ease: 'Power1',
        })
      } else {
        // Attack — flash and chance to kill
        this.cameras.main.shake(50, 0.002)
        if (Math.random() < 0.35) {
          this.killUnit(target)
        }
      }
    }

    // Defenders fight back
    for (const defender of aliveDefenders) {
      const target = this.findNearest(defender, aliveAttackers)
      if (!target) continue
      const dist = Phaser.Math.Distance.Between(defender.x, defender.y, target.x, target.y)
      if (dist <= 50 && Math.random() < 0.3) {
        this.killUnit(target)
      }
    }
  }

  findNearest(unit: Phaser.GameObjects.Image, targets: Phaser.GameObjects.Image[]): Phaser.GameObjects.Image | null {
    let nearest: Phaser.GameObjects.Image | null = null
    let minDist = Infinity
    for (const t of targets) {
      const d = Phaser.Math.Distance.Between(unit.x, unit.y, t.x, t.y)
      if (d < minDist) { minDist = d; nearest = t }
    }
    return nearest
  }

  killUnit(unit: Phaser.GameObjects.Image) {
    // Explosion effect
    const explosion = this.add.circle(unit.x, unit.y, 12, 0xff6600, 0.8).setDepth(60)
    this.tweens.add({ targets: explosion, scale: 2, alpha: 0, duration: 400, onComplete: () => explosion.destroy() })
    unit.destroy()
  }
}
