package main

import (
	"bufio"
	"os"
)

type template struct {
	typeList     []string
	power        int
	toughness    int
	displayGroup string
}

type vanilla struct {
	*template
	amount  int
	damage  int
	dying   bool
	arcbond bool
	control bool
	vip     bool
	blast   bool
	swap    bool
}
type bishop struct {
	*vanilla
	triggerOn       string
	create          *template
	waitingTriggers int
	loud            bool
}

type gamestate struct {
	bishopList      []*bishop
	vanillaList     []*vanilla
	deadBishops     []*bishop
	coatNumber      int
	arcbondDamage   int
	arcbondsOnStack int
	opponentLife    int
	loud            bool
	vipDeath        bool
	displayGroups   []string
}

func setCommonValues(v *vanilla, pl *parsedLine) {
	v.arcbond = listContains(pl.keywords, "Arcbond")
	v.control = listContains(pl.keywords, "Control")
	v.vip = listContains(pl.keywords, "VIP")
	v.blast = listContains(pl.keywords, "Blast")
	v.swap = listContains(pl.keywords, "Swap")
	v.amount = pl.multiplicity
	v.displayGroup = pl.displayGroup1
}

func (gs *gamestate) addVanilla(pl *parsedLine) {
	newVanilla := vanilla{
		template: &template{
			typeList:  []string{"Golem"},
			power:     3,
			toughness: 3,
		},
	}
	setCommonValues(&newVanilla, pl)
	if value, set := pl.typeReplacements[1]; set {
		newVanilla.typeList[0] = value
	}
	gs.vanillaList = append(gs.vanillaList, &newVanilla)
}

func (gs *gamestate) addBishop(pl *parsedLine) {
	newBishop := bishop{
		vanilla: &vanilla{
			template: &template{
				typeList:  []string{"Human", "Cleric"},
				power:     1,
				toughness: 4,
			},
		},
		triggerOn: "Angel",
		create: &template{
			typeList:     []string{"Spirit"},
			power:        1,
			toughness:    1,
			displayGroup: pl.displayGroup2,
		},
		loud: listContains(pl.keywords, "Loud"),
	}
	setCommonValues(newBishop.vanilla, pl)
	if value, set := pl.typeReplacements[1]; set {
		newBishop.typeList[0] = value
	}
	if value, set := pl.typeReplacements[2]; set {
		newBishop.typeList[1] = value
	}
	if newBishop.typeList[0] == newBishop.typeList[1] {
		newBishop.typeList = newBishop.typeList[:1]
	}

	if value, set := pl.typeReplacements[3]; set {
		newBishop.triggerOn = value
	}
	if value, set := pl.typeReplacements[4]; set {
		newBishop.create.typeList[0] = value
	}
	gs.bishopList = append(gs.bishopList, &newBishop)

}

func (t *template) applyDralnu(pl *parsedLine) {
	type1 := "Goblin"
	if value, set := pl.typeReplacements[1]; set {
		type1 = value
	}
	type2 := "Zombie"
	if value, set := pl.typeReplacements[2]; set {
		type2 = value
	}

	if listContains(t.typeList, type1) {
		t.power += pl.multiplicity
		t.toughness += pl.multiplicity
		if !listContains(t.typeList, type2) {
			t.typeList = append(t.typeList, type2)
		}
	}
}

func (b *bishop) applyDralnu(pl *parsedLine) {
	b.template.applyDralnu(pl)
	b.create.applyDralnu(pl)
}

func (gs *gamestate) adjustBaselinesForDralnu(pl *parsedLine) {
	for _, bishop := range gs.bishopList {
		bishop.applyDralnu(pl)
	}
	for _, vanilla := range gs.vanillaList {
		vanilla.template.applyDralnu(pl)
	}
}

func (t *template) applyCoats(gs *gamestate) {
	counter := -1
	for _, bishop := range gs.bishopList {
		if listsIntersect(t.typeList, bishop.typeList) {
			counter += bishop.amount
		}
	}
	for _, vanilla := range gs.vanillaList {
		if listsIntersect(t.typeList, vanilla.typeList) {
			counter += vanilla.amount
		}
	}
	t.power += gs.coatNumber * counter
	t.toughness += gs.coatNumber * counter
}

func (gs *gamestate) applyCoats() {
	for _, bishop := range gs.bishopList {
		bishop.template.applyCoats(gs)
	}
	for _, vanilla := range gs.vanillaList {
		vanilla.template.applyCoats(gs)
	}
}

func buildInitialGamestate() *gamestate {
	gs := gamestate{
		bishopList:      []*bishop{},
		vanillaList:     []*vanilla{},
		deadBishops:     []*bishop{},
		coatNumber:      0,
		arcbondDamage:   0,
		arcbondsOnStack: 0,
		opponentLife:    20,
		loud:            false,
	}
	dralnus := []*parsedLine{}

	programByLine := bufio.NewScanner(os.Stdin)
	for programByLine.Scan() {

		addedObject := parseProgramLine(programByLine.Text())
		if addedObject == nil {
			//ignore invalid lines
			continue
		}

		if addedObject.displayGroup1 != "" && !listContains(gs.displayGroups, addedObject.displayGroup1) {
			gs.displayGroups = append(gs.displayGroups, addedObject.displayGroup1)
		}
		if addedObject.displayGroup2 != "" && !listContains(gs.displayGroups, addedObject.displayGroup2) {
			gs.displayGroups = append(gs.displayGroups, addedObject.displayGroup2)
		}

		if addedObject.multiplicity <= 0 {
			continue
		}
		switch addedObject.objectType {
		case "Coat":
			gs.coatNumber += addedObject.multiplicity
		case "Vanilla":
			gs.addVanilla(addedObject)
		case "Bishop":
			gs.addBishop(addedObject)
		case "Dralnu":
			dralnus = append(dralnus, addedObject)
		}

	}

	//We apply all Dralnu's in the order that they appear in the program.
	//We ignore dependencies and layers.
	//Please don't build complicated Dralnu chains.
	for _, dralnu := range dralnus {
		gs.adjustBaselinesForDralnu(dralnu)
	}
	gs.applyCoats()

	return &gs
}
