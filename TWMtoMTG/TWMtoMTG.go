package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type programLine struct {
	quantity    int
	object      string
	typeChanges []string
	keywords    []string
	display1    string
	display2    string
}

func (pl programLine) String() string {
	result := fmt.Sprintf("%dx %s", pl.quantity, pl.object)
	for key, value := range pl.typeChanges {
		result += fmt.Sprintf(" %d:%s", key+1, value)
	}
	for _, value := range pl.keywords {
		result += fmt.Sprintf(" %s", value)
	}
	if pl.display1 != "" {
		result += fmt.Sprintf(" <%s", pl.display1)
	}
	if pl.display2 != "" {
		result += fmt.Sprintf(" >%s", pl.display2)
	}
	return result
}

//list of creature creatureTypes
//no Goblin, Zombie, Human, Golem, Wall
var creatureTypes = []string{"Advisor", "Aetherborn", "Alien", "Ally", "Angel", "Antelope", "Ape", "Archer", "Archon", "Army", "Artificer", "Assassin", "Assembly-Worker", "Astartes", "Atog", "Aurochs", "Avatar", "Azra", "Badger", "Balloon", "Barbarian", "Bard", "Basilisk", "Bat", "Bear", "Beast", "Beeble", "Beholder", "Berserker", "Bird", "Blinkmoth", "Boar", "Bringer", "Brushwagg", "Camarid", "Camel", "Caribou", "Carrier", "Cat", "Centaur", "Cephalid", "Child", "Chimera", "Citizen", "Cleric", "Clown", "Cockatrice", "Construct", "Coward", "Crab", "Crocodile", "C’tan", "Custodes", "Cyclops", "Dauthi", "Demigod", "Demon", "Deserter", "Devil", "Dinosaur", "Djinn", "Dog", "Dragon", "Drake", "Dreadnought", "Drone", "Druid", "Dryad", "Dwarf", "Efreet", "Egg", "Elder", "Eldrazi", "Elemental", "Elephant", "Elf", "Elk", "Employee", "Eye", "Faerie", "Ferret", "Fish", "Flagbearer", "Fox", "Fractal", "Frog", "Fungus", "Gamer", "Gargoyle", "Germ", "Giant", "Gith", "Gnoll", "Gnome", "Goat", "God", "Gorgon", "Graveborn", "Gremlin", "Griffin", "Guest", "Hag", "Halfling", "Hamster", "Harpy", "Hellion", "Hippo", "Hippogriff", "Homarid", "Homunculus", "Horror", "Horse", "Hydra", "Hyena", "Illusion", "Imp", "Incarnation", "Inkling", "Inquisitor", "Insect", "Jackal", "Jellyfish", "Juggernaut", "Kavu", "Kirin", "Kithkin", "Knight", "Kobold", "Kor", "Kraken", "Lamia", "Lammasu", "Leech", "Leviathan", "Lhurgoyf", "Licid", "Lizard", "Manticore", "Masticore", "Mercenary", "Merfolk", "Metathran", "Minion", "Minotaur", "Mite", "Mole", "Monger", "Mongoose", "Monk", "Monkey", "Moonfolk", "Mouse", "Mutant", "Myr", "Mystic", "Naga", "Nautilus", "Necron", "Nephilim", "Nightmare", "Nightstalker", "Ninja", "Noble", "Noggle", "Nomad", "Nymph", "Octopus", "Ogre", "Ooze", "Orb", "Orc", "Orgg", "Otter", "Ouphe", "Ox", "Oyster", "Pangolin", "Peasant", "Pegasus", "Pentavite", "Performer", "Pest", "Phelddagrif", "Phoenix", "Phyrexian", "Pilot", "Pincher", "Pirate", "Plant", "Praetor", "Primarch", "Prism", "Processor", "Rabbit", "Raccoon", "Ranger", "Rat", "Rebel", "Reflection", "Rhino", "Rigger", "Robot", "Rogue", "Sable", "Salamander", "Samurai", "Sand", "Saproling", "Satyr", "Scarecrow", "Scion", "Scorpion", "Scout", "Sculpture", "Serf", "Serpent", "Servo", "Shade", "Shaman", "Shapeshifter", "Shark", "Sheep", "Siren", "Skeleton", "Slith", "Sliver", "Slug", "Snake", "Soldier", "Soltari", "Spawn", "Specter", "Spellshaper", "Sphinx", "Spider", "Spike", "Spirit", "Splinter", "Sponge", "Squid", "Squirrel", "Starfish", "Surrakar", "Survivor", "Tentacle", "Tetravite", "Thalakos", "Thopter", "Thrull", "Tiefling", "Treefolk", "Trilobite", "Triskelavite", "Troll", "Turtle", "Tyranid", "Unicorn", "Vampire", "Vedalken", "Viashino", "Volver", "Walrus", "Warlock", "Warrior", "Weird", "Werewolf", "Whale", "Wizard", "Wolf", "Wolverine", "Wombat", "Worm", "Wraith", "Wurm", "Yeti", "Zubera"}

//returns the creature type for clock_i. j=0 for the special type, j=1 for the support type.
//will panic if it can't find a type!
func clockType(i, j int) string {
	switch i {
	case 1:
		if j == 0 {
			return "Goblin"
		}
		return "Zombie"
	case 2:
		if j == 0 {
			return "Human"
		}
		return "Golem"
	default:
		return creatureTypes[-2+2*i+j]
	}
}

func invalidInput() {
	fmt.Fprintln(os.Stdout, "Invalid Program")
	os.Exit(1)
}

func main() {
	//interpret input
	scanner := bufio.NewScanner(os.Stdin)
	programText := ""
	for scanner.Scan() {
		programText += strings.ReplaceAll(scanner.Text(), " ", "")
	}
	programText = strings.Trim(programText, "[]")
	programLines := strings.Split(programText, "],[")
	numClocks := len(programLines) - 1
	if numClocks < 1 {
	}
	row0Texts := strings.Split(programLines[0], ",")
	if len(row0Texts) != numClocks+1 {
		invalidInput()
	}
	max, err := strconv.Atoi(row0Texts[0])
	if err != nil || max < numClocks {
		invalidInput()

	}
	row0 := []int{max}
	for _, numberString := range row0Texts[1:] {
		number, err := strconv.Atoi(numberString)
		if err != nil || number != numClocks {
			invalidInput()
		}
		row0 = append(row0, number)
	}
	programMatrix := [][]int{row0}
	for _, newLine := range programLines[1:] {
		numberTexts := strings.Split(newLine, ",")
		if len(numberTexts) != numClocks+1 {
			invalidInput()
		}
		newMatrixRow := []int{}
		for _, numberString := range numberTexts {
			number, err := strconv.Atoi(numberString)
			if err != nil || number > max {
				invalidInput()
			}
			newMatrixRow = append(newMatrixRow, number)
		}
		programMatrix = append(programMatrix, newMatrixRow)
	}
	//build TWM program
	if numClocks > 138 {
		fmt.Fprintln(os.Stdout, "Too many clocks")
		os.Exit(1)
	}
	program := []any{}
	program = append(program,
		"--setup",
		programLine{
			quantity: 3,
			object:   "Coat",
		},
		programLine{
			quantity:    1,
			object:      "Vanilla",
			typeChanges: []string{creatureTypes[1]},
			keywords:    []string{"Control"},
		},
		programLine{
			quantity: 1,
			object:   "Vanilla",
			keywords: []string{"Swap"},
		},
		"",
		"--heartbeat",
		programLine{
			quantity:    1,
			object:      "Vanilla",
			typeChanges: []string{creatureTypes[3]},
		},
		programLine{
			quantity:    1,
			object:      "Bishop",
			typeChanges: []string{creatureTypes[0], creatureTypes[2], creatureTypes[3], creatureTypes[3]},
		},
		programLine{
			quantity:    2,
			object:      "Bishop",
			typeChanges: []string{creatureTypes[0], creatureTypes[0], creatureTypes[3], creatureTypes[0]},
		},
		"",
		"--halt",
		programLine{
			quantity:    1,
			object:      "Vanilla",
			typeChanges: []string{creatureTypes[2]},
			keywords:    []string{"Arcbond", "Blast"},
		},
		programLine{
			quantity:    1,
			object:      "Bishop",
			typeChanges: []string{creatureTypes[0], creatureTypes[2], creatureTypes[3], creatureTypes[2]},
		})
	for i, clock := range programMatrix {
		if i == 0 {
			continue
		}
		program = append(program,
			"",
			fmt.Sprintf("--clock %d", i),
			programLine{
				quantity:    1,
				object:      "Dralnu",
				typeChanges: []string{clockType(i, 0), clockType(i, 1)},
			},
			programLine{
				quantity:    1,
				object:      "Vanilla",
				typeChanges: []string{clockType(i, 0)},
				display1:    fmt.Sprintf("clock_%d", i),
			},
			programLine{
				quantity:    clock[0] + 1,
				object:      "Vanilla",
				typeChanges: []string{clockType(i, 1)},
				display1:    fmt.Sprintf("clock_%d", i),
			})
		if clock[i] <= 0 {
			//halting clock
			program = append(program,
				programLine{
					quantity:    1,
					object:      "Bishop",
					typeChanges: []string{creatureTypes[0], creatureTypes[0], clockType(i, 0), creatureTypes[3]},
				},
				"--output preservation")
			for j := range clock {
				if j != 0 && j != i {
					program = append(program,
						programLine{
							quantity:    3,
							object:      "Bishop",
							typeChanges: []string{creatureTypes[0], creatureTypes[0], clockType(i, 0), clockType(j, 1)},
							display2:    fmt.Sprintf("clock_%d", j),
						})
				}
			}
			continue
		}
		//normal clock
		program = append(program,
			programLine{
				quantity:    1,
				object:      "Bishop",
				typeChanges: []string{creatureTypes[0], creatureTypes[0], clockType(i, 0), clockType(i, 0)},
				display2:    fmt.Sprintf("clock_%d", i),
			})
		for j, q := range clock {
			if j == 0 {
				continue
			}
			if j == i {
				q -= 1
			}
			program = append(program,
				programLine{
					quantity:    q,
					object:      "Bishop",
					typeChanges: []string{creatureTypes[0], creatureTypes[0], clockType(i, 0), clockType(j, 1)},
					display2:    fmt.Sprintf("clock_%d", j),
				})
		}
	}
	//print TWM program
	for _, line := range program {
		fmt.Println(line)
	}
}
