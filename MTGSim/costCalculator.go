package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/rgeoghegan/tabulate"
)

type costItem struct {
	text          string
	unitPrice     int
	golemDiscount bool
	amount        int
}

func (ci costItem) strUnitPrice() string {
	if ci.golemDiscount {
		return "1/3"
	}
	return strconv.Itoa(ci.unitPrice)
}

func (ci costItem) strAmount() string {
	return strconv.Itoa(ci.amount)
}

func (ci costItem) total() int {
	if ci.golemDiscount {
		return (ci.amount + 2) / 3 //+2 because integer division gives the floor()
	}
	return ci.amount * ci.unitPrice
}

func (ci costItem) strTotal() string {
	return strconv.Itoa(ci.total())
}

func writeReceipt(cl []costItem) {
	total := 0
	headers := []string{"Description", "Quantity", "Price/Unit", "LineTotal"}
	table := [][]string{}
	for _, item := range cl {
		table = append(table, []string{item.text, item.strAmount(), item.strUnitPrice(), item.strTotal()})
		total += item.total()
	}
	table = append(table, []string{"", "", "", ""})
	table = append(table, []string{"", "", "Total:", strconv.Itoa(total)})
	receiptText, err := tabulate.Tabulate(table, &tabulate.Layout{Headers: headers, Format: tabulate.FancyGridFormat})
	if err != nil {
		panic(err)
	}
	fmt.Println(receiptText)
}

func calculateCost() {
	programByLine := bufio.NewScanner(os.Stdin)

	coatCost := costItem{
		text:      "Coat of Arms",
		unitPrice: 1,
		amount:    0,
	}
	vanillaCost := costItem{
		text:          "Vanilla",
		golemDiscount: true,
		amount:        -6,
	}
	bishopCost := costItem{
		text:      "Bishop of Wings",
		unitPrice: 1,
		amount:    0,
	}
	dralnuCost := costItem{
		text:      "Dralnu's Crusade",
		unitPrice: 1,
		amount:    0,
	}
	dralnuChangeCost := costItem{
		text:      "Artificial Evolution on Dralnu's Crusade",
		unitPrice: 2,
		amount:    0,
	}

	dralnuSwapCost := costItem{
		text:      "Extra AE for Zombie->Goblin Dralnu's",
		unitPrice: 1,
		amount:    0,
	}

	cheapArtificialEvolution := false

	for programByLine.Scan() {

		addedObject := parseProgramLine(programByLine.Text())
		if addedObject == nil {
			//ignore invalid lines
			continue
		}

		switch addedObject.objectType {
		case "Coat":
			coatCost.amount += addedObject.multiplicity
		case "Vanilla":
			vanillaCost.amount += addedObject.multiplicity
		case "Bishop":
			bishopCost.amount += addedObject.multiplicity
		case "Dralnu":
			dralnuCost.amount += addedObject.multiplicity
			dralnuChangeCost.amount += len(addedObject.typeReplacements) * addedObject.multiplicity
			type1 := addedObject.typeReplacements[1]
			type2 := addedObject.typeReplacements[2]
			if type1 == "Goblin" {
				dralnuChangeCost.amount -= addedObject.multiplicity
			}
			if type2 == "Zombie" {
				dralnuChangeCost.amount -= addedObject.multiplicity
			}
			if type1 == "Zombie" && type2 == "Goblin" {
				dralnuSwapCost.amount += addedObject.multiplicity
			}
			if (type1 == "Human" || type1 == "Cleric") && type2 == "Golem" && addedObject.multiplicity > 0 {
				cheapArtificialEvolution = true
			}
		}

	}
	cl := []costItem{
		{
			text:      "Replication Technique recovery",
			unitPrice: 1,
			amount:    1,
		},
		coatCost,
		{
			text:      "Free Vanilla (from start up)",
			unitPrice: 0,
			amount:    6,
		},
		vanillaCost,
	}
	if vanillaCost.amount%3 != 0 {
		cl = append(cl,
			costItem{
				text:      "Missing out on free vanilla!",
				unitPrice: 0,
				amount:    3 - (vanillaCost.amount % 3),
			})
	}
	cl = append(cl,
		bishopCost,
		dralnuCost,
		dralnuChangeCost,
	)
	if dralnuSwapCost.amount > 0 {
		cl = append(cl, dralnuSwapCost)
	}
	cl = append(cl,
		costItem{
			text:      "Artificial Evolution on all creatures",
			unitPrice: 1,
			amount:    1,
		})
	if !cheapArtificialEvolution {
		cl = append(cl,
			costItem{
				text:      "Extra AE without Human/Cleric->Golem Dralnu's",
				unitPrice: 2,
				amount:    1,
			})
	}
	cl = append(cl,
		costItem{
			text:      "Arcbond (from hand)",
			unitPrice: 0,
			amount:    1,
		},
		costItem{
			text:      "Comeuppance",
			unitPrice: 1,
			amount:    1,
		},
		costItem{
			text:      "Scrambleverse",
			unitPrice: 1,
			amount:    1,
		},
		costItem{
			text:      "Soulblast (computation start)",
			unitPrice: 1,
			amount:    1,
		},
		costItem{
			text:      "Scrambleverse (after computation)",
			unitPrice: 2,
			amount:    1,
		},
		costItem{
			text:      "Soulblast (after computation)",
			unitPrice: 2,
			amount:    1,
		},
	)
	writeReceipt(cl)
}
