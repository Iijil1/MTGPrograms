package main

import (
	"fmt"
	"os"
)

func initialSoulblast(gs *gamestate) {
	resolveSwap(gs)
	castSoulblast(gs)
	resolveSoulblast(gs)
}

func resolveSwap(gs *gamestate) {
	counter := 0
	for _, bishop := range gs.bishopList {
		if bishop.swap {
			bishop.dying = true
			counter += bishop.amount
		}
	}
	for _, vanilla := range gs.vanillaList {
		if vanilla.swap {
			vanilla.dying = true
			counter += vanilla.amount
		}
	}
	if counter == 0 {
		fmt.Fprintln(os.Stderr, "No Audacious Swap target")
		os.Exit(1)
	}
	if counter > 1 {
		fmt.Fprintln(os.Stderr, "Too many Audacious Swap targets")
		os.Exit(1)
	}
	removeDead(gs)
}

func castSoulblast(gs *gamestate) {
	for _, bishop := range gs.bishopList {
		if bishop.control {
			bishop.dying = true
			triggerBishops(gs, bishop.vanilla)
			gs.arcbondDamage += bishop.power * bishop.amount
		}
	}
	for _, vanilla := range gs.vanillaList {
		if vanilla.control {
			vanilla.dying = true
			triggerBishops(gs, vanilla)
			gs.arcbondDamage += vanilla.power * vanilla.amount
		}
	}
	removeDead(gs)
	resolveBishopTriggers(gs)
}

func resolveSoulblast(gs *gamestate) {
	counter := 0
	for _, bishop := range gs.bishopList {
		if bishop.blast {
			bishop.damage += gs.arcbondDamage
			counter += bishop.amount
			if bishop.arcbond {
				gs.arcbondsOnStack += bishop.amount
			}
		}
	}
	for _, vanilla := range gs.vanillaList {
		if vanilla.blast {
			vanilla.damage += gs.arcbondDamage
			counter += vanilla.amount
			if vanilla.arcbond {
				gs.arcbondsOnStack += vanilla.amount
			}
		}
	}
	if counter == 0 {
		fmt.Fprintln(os.Stderr, "No Soulblast target")
		os.Exit(1)
	}
	if counter > 1 {
		fmt.Fprintln(os.Stderr, "Too many Soulblast targets")
		os.Exit(1)
	}
	stateBasedActions(gs)
	resolveBishopTriggers(gs)
}

func addCoatEffect(gs *gamestate, v *vanilla) {
	for _, bishop := range gs.bishopList {
		if listsIntersect(v.typeList, bishop.typeList) {
			bishop.power += gs.coatNumber * v.amount
			bishop.toughness += gs.coatNumber * v.amount
			v.power += gs.coatNumber * bishop.amount
			v.toughness += gs.coatNumber * bishop.amount
		}
		if listContains(v.typeList, bishop.triggerOn) && !v.control && !bishop.control {
			gs.opponentLife += 4 * v.amount * bishop.amount
		}
	}
	for _, vanilla := range gs.vanillaList {
		if listsIntersect(v.typeList, vanilla.typeList) {
			vanilla.power += gs.coatNumber * v.amount
			vanilla.toughness += gs.coatNumber * v.amount
			v.power += gs.coatNumber * vanilla.amount
			v.toughness += gs.coatNumber * vanilla.amount
		}
	}
	v.power += gs.coatNumber * (v.amount - 1)
	v.toughness += gs.coatNumber * (v.amount - 1)
}

func resolveBishopTrigger(gs *gamestate, b *bishop) {
	template := *b.create
	newVanilla := &vanilla{
		template: &template,
		amount:   b.amount * b.waitingTriggers,
		control:  b.control,
	}
	addCoatEffect(gs, newVanilla)
	gs.vanillaList = append(gs.vanillaList, newVanilla)
	b.waitingTriggers = 0
	if b.loud {
		gs.loud = true
	}
}

func resolveBishopTriggers(gs *gamestate) {
	for _, bishop := range gs.bishopList {
		if bishop.waitingTriggers > 0 {
			resolveBishopTrigger(gs, bishop)
		}
	}
	for _, bishop := range gs.deadBishops {
		resolveBishopTrigger(gs, bishop)
	}
	gs.deadBishops = gs.deadBishops[:0]
}

func triggerBishops(gs *gamestate, v *vanilla) {
	for _, bishop := range gs.bishopList {
		if bishop.control == v.control && listContains(v.typeList, bishop.triggerOn) {
			bishop.waitingTriggers += v.amount
		}
	}
	if v.vip {
		gs.loud = true
		gs.vipDeath = true
	}
}

func checkDying(gs *gamestate) bool {
	somethingDies := false
	for _, bishop := range gs.bishopList {
		if bishop.damage >= bishop.toughness {
			bishop.dying = true
			somethingDies = true
			triggerBishops(gs, bishop.vanilla)
		}
	}
	for _, vanilla := range gs.vanillaList {
		if vanilla.damage >= vanilla.toughness {
			vanilla.dying = true
			somethingDies = true
			triggerBishops(gs, vanilla)
		}
	}
	return somethingDies
}

func removeCoatEffect(gs *gamestate, v *vanilla) {
	for _, bishop := range gs.bishopList {
		if listsIntersect(v.typeList, bishop.typeList) {
			bishop.power -= gs.coatNumber * v.amount
			bishop.toughness -= gs.coatNumber * v.amount
		}
	}
	for _, vanilla := range gs.vanillaList {
		if listsIntersect(v.typeList, vanilla.typeList) {
			vanilla.power -= gs.coatNumber * v.amount
			vanilla.toughness -= gs.coatNumber * v.amount
		}
	}
}

func removeDead(gs *gamestate) {
	for i := len(gs.bishopList) - 1; i >= 0; i -= 1 {
		bishop := gs.bishopList[i]
		if bishop.dying {
			if bishop.waitingTriggers > 0 {
				gs.deadBishops = append(gs.deadBishops, bishop)
			}
			removeCoatEffect(gs, bishop.vanilla)
			removeFromList(&gs.bishopList, i)
		}
	}
	for i := len(gs.vanillaList) - 1; i >= 0; i -= 1 {
		vanilla := gs.vanillaList[i]
		if vanilla.dying {
			removeCoatEffect(gs, vanilla)
			removeFromList(&gs.vanillaList, i)
		}
	}
}

func stateBasedActions(gs *gamestate) {
	for checkDying(gs) {
		removeDead(gs)
	}
}

func resolveArcbond(gs *gamestate) {
	gs.arcbondsOnStack -= 1
	gs.opponentLife -= gs.arcbondDamage
	for _, bishop := range gs.bishopList {
		bishop.damage += gs.arcbondDamage
		if bishop.arcbond {
			gs.arcbondsOnStack += bishop.amount
		}
	}
	for _, vanilla := range gs.vanillaList {
		vanilla.damage += gs.arcbondDamage
		if vanilla.arcbond {
			gs.arcbondsOnStack += vanilla.amount
		}
	}
}

func simulateStep(gs *gamestate) {
	gs.loud = false
	resolveArcbond(gs)
	stateBasedActions(gs)
	resolveBishopTriggers(gs)
}

func simulateProgram(limit int) {

	gs := buildInitialGamestate()
	initialSoulblast(gs)
	displayState(gs)
	steps := 1
	for gs.arcbondsOnStack > 0 {

		simulateStep(gs)
		if gs.loud {
			displayState(gs)
		}
		if gs.vipDeath {
			fmt.Println("A VIP died, we will not profit from this computation!")
			break
		}
		if gs.opponentLife <= 0 {
			fmt.Println("The opponent died, we will not profit from this computation!")
			break
		}
		steps += 1
		if steps > limit {
			fmt.Println("Oops, thats too much Arcbond! We need to stop while we can.")
			break
		}
	}
	fmt.Println("Final gamestate:")
	displayState(gs)
}
