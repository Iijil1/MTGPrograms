package main

import (
	"fmt"
	"strconv"

	"github.com/rgeoghegan/tabulate"
)

func (gs *gamestate) String() string {
	return fmt.Sprintf("%+v", *gs)
}

func (b bishop) String() string {
	return fmt.Sprintf("%dx %d/%d -%d Bishop of Wings: %v, %v -> %d/%d %v", b.amount, b.power, b.toughness, b.damage, b.typeList, b.triggerOn, b.create.power, b.create.toughness, b.create.typeList)
}

func (v vanilla) String() string {
	return fmt.Sprintf("%dx %d/%d -%d Vanilla: %v", v.amount, v.power, v.toughness, v.damage, v.typeList)
}

func displayFullState(gs *gamestate) {
	fmt.Println("Opponents Life:", gs.opponentLife)
	fmt.Println("Arcbond triggers:", gs.arcbondsOnStack)
	fmt.Println("Damage per trigger:", gs.arcbondDamage)
	fmt.Printf("%dx Coat of Arms\n", gs.coatNumber)
	for _, bishop := range gs.bishopList {
		fmt.Println(bishop)
	}
	for _, vanilla := range gs.vanillaList {
		fmt.Println(vanilla)
	}
}

type displayGroup struct {
	size int
	life int
}

func displayState(gs *gamestate) {
	groups := map[string]*displayGroup{}
	for _, name := range gs.displayGroups {
		groups[name] = &displayGroup{0, 0}
	}
	for _, bishop := range gs.bishopList {
		group := groups[bishop.displayGroup]
		if group == nil {
			continue
		}
		group.size += bishop.amount
		if group.life > bishop.toughness-bishop.damage || group.life == 0 {
			group.life = bishop.toughness - bishop.damage
		}
	}
	for _, vanilla := range gs.vanillaList {
		group := groups[vanilla.displayGroup]
		if group == nil {
			continue
		}
		group.size += vanilla.amount
		if group.life > vanilla.toughness-vanilla.damage || group.life == 0 {
			group.life = vanilla.toughness - vanilla.damage
		}
	}

	headers := []string{"Group", "Size", "Life"}
	table := [][]string{}
	for _, name := range gs.displayGroups {
		group := groups[name]
		if group == nil {
			continue
		}
		table = append(table, []string{name, strconv.Itoa(group.size), strconv.Itoa(group.life)})
	}
	battlefieldText, err := tabulate.Tabulate(table, &tabulate.Layout{Headers: headers, Format: tabulate.FancyGridFormat})
	if err != nil {
		panic(err)
	}
	fmt.Println(battlefieldText)
}
