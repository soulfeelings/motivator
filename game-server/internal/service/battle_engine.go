package service

import (
	"context"
	"math/rand"

	"github.com/hustlers/motivator-game/internal/model"
)

type combatant struct {
	unitID  string
	hp      int
	attack  int
	defense int
	side    string
}

// Attack runs auto-battle between attacker and defender bases.
func (s *GameService) Attack(ctx context.Context, attackerBaseID, defenderBaseID string) (*model.Battle, error) {
	if attackerBaseID == defenderBaseID {
		return nil, model.ErrCannotAttackSelf
	}

	attackerArmy, err := s.army.GetArmy(ctx, attackerBaseID)
	if err != nil {
		return nil, err
	}
	defenderArmy, err := s.army.GetArmy(ctx, defenderBaseID)
	if err != nil {
		return nil, err
	}

	unitTypes, err := s.army.GetUnitTypes(ctx)
	if err != nil {
		return nil, err
	}
	stats := make(map[string]model.UnitType)
	for _, ut := range unitTypes {
		stats[ut.ID] = ut
	}

	var attackers, defenders []combatant
	for _, au := range attackerArmy {
		ut := stats[au.UnitID]
		for range au.Count {
			attackers = append(attackers, combatant{unitID: au.UnitID, hp: ut.HP, attack: ut.Attack, defense: ut.Defense, side: "attacker"})
		}
	}
	for _, du := range defenderArmy {
		ut := stats[du.UnitID]
		for range du.Count {
			defenders = append(defenders, combatant{unitID: du.UnitID, hp: ut.HP, attack: ut.Attack, defense: ut.Defense, side: "defender"})
		}
	}

	// Turret bonus for defender
	defBuildings, _ := s.bases.ListBuildings(ctx, defenderBaseID)
	turretBonus := 0
	for _, b := range defBuildings {
		if b.BuildingID == "turret" {
			turretBonus += 20 * b.Level
		}
	}

	// Simulate
	var replay []model.ReplayFrame
	for tick := 0; tick < 30 && countAlive(attackers) > 0 && countAlive(defenders) > 0; tick++ {
		var events []model.ReplayEvent

		for i := range attackers {
			if attackers[i].hp <= 0 {
				continue
			}
			t := pickAliveIdx(defenders)
			if t < 0 {
				break
			}
			dmg := calcDamage(attackers[i].attack, defenders[t].defense)
			defenders[t].hp -= dmg
			events = append(events, model.ReplayEvent{Type: "attack", UnitID: attackers[i].unitID, TargetID: defenders[t].unitID, Damage: dmg, Side: "attacker"})
		}

		bonusPerUnit := 0
		if alive := countAlive(defenders); alive > 0 {
			bonusPerUnit = turretBonus / alive
		}
		for i := range defenders {
			if defenders[i].hp <= 0 {
				continue
			}
			t := pickAliveIdx(attackers)
			if t < 0 {
				break
			}
			dmg := calcDamage(defenders[i].attack+bonusPerUnit, attackers[t].defense)
			attackers[t].hp -= dmg
			events = append(events, model.ReplayEvent{Type: "attack", UnitID: defenders[i].unitID, TargetID: attackers[t].unitID, Damage: dmg, Side: "defender"})
		}

		// Death events
		for i := range attackers {
			if attackers[i].hp <= 0 && attackers[i].hp > -9999 {
				events = append(events, model.ReplayEvent{Type: "death", UnitID: attackers[i].unitID, Side: "attacker"})
				attackers[i].hp = -10000 // mark processed
			}
		}
		for i := range defenders {
			if defenders[i].hp <= 0 && defenders[i].hp > -9999 {
				events = append(events, model.ReplayEvent{Type: "death", UnitID: defenders[i].unitID, Side: "defender"})
				defenders[i].hp = -10000
			}
		}

		replay = append(replay, model.ReplayFrame{Tick: tick, Events: events})
	}

	attackerSurvived := countAlive(attackers)
	defenderSurvived := countAlive(defenders)
	attackerWon := attackerSurvived > defenderSurvived

	var winnerID *string
	coinsWon := 50
	xpWon := 100
	if attackerWon {
		winnerID = &attackerBaseID
	} else {
		winnerID = &defenderBaseID
	}

	attackerLost := calcUnitLosses(attackerArmy, attackers)
	defenderLost := calcUnitLosses(defenderArmy, defenders)

	for unitID, lost := range attackerLost {
		s.army.RemoveUnits(ctx, attackerBaseID, unitID, lost)
	}
	for unitID, lost := range defenderLost {
		s.army.RemoveUnits(ctx, defenderBaseID, unitID, lost)
	}

	if winnerID != nil {
		s.bases.UpdateCoins(ctx, *winnerID, coinsWon)
	}

	battle := &model.Battle{
		AttackerID:   attackerBaseID,
		DefenderID:   defenderBaseID,
		WinnerID:     winnerID,
		AttackerLost: attackerLost,
		DefenderLost: defenderLost,
		ReplayData:   replay,
		CoinsWon:     coinsWon,
		XPWon:        xpWon,
	}

	if err := s.battles.Create(ctx, battle); err != nil {
		return nil, err
	}

	return battle, nil
}

func countAlive(units []combatant) int {
	n := 0
	for _, u := range units {
		if u.hp > 0 {
			n++
		}
	}
	return n
}

func pickAliveIdx(units []combatant) int {
	alive := make([]int, 0)
	for i, u := range units {
		if u.hp > 0 {
			alive = append(alive, i)
		}
	}
	if len(alive) == 0 {
		return -1
	}
	return alive[rand.Intn(len(alive))]
}

func calcDamage(attack, defense int) int {
	dmg := attack - defense/2 + rand.Intn(5)
	if dmg < 1 {
		dmg = 1
	}
	return dmg
}

func calcUnitLosses(original []model.ArmyUnit, combatants []combatant) map[string]int {
	survived := make(map[string]int)
	for _, c := range combatants {
		if c.hp > 0 {
			survived[c.unitID]++
		}
	}
	losses := make(map[string]int)
	for _, au := range original {
		lost := au.Count - survived[au.UnitID]
		if lost > 0 {
			losses[au.UnitID] = lost
		}
	}
	return losses
}
