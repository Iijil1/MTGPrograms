A simulator to run the mandatory actions in the very specific MTG boardstate that is set up by the champion 14 card deck.

# Flags:<br>
-c<br>
-- Calculate the cost of the boardstate in terms of Audacious Swap casts required to set it up, instead of simulating it.<br>
-l<br>
-- Set the upper limit of Arcbond triggers that are resolved before giving up on the simulation. Defaults to 5000

Input the boardstate via stdin.

# Input format:

Every line describes an object on the board.

Example Lines:

1x Coat<br>
-- 1 Coat of Arms is on the battlefield

4x Vanilla 1:Ally<br>
-- there are 4 3/3 creatures on the battlefield under the opponents control. They have the creature type Ally and no abilities relevant for the simulation. Those will be either Precursor Golems or the Golem tokens created by it, modified by Artificial Evolution.

2x Bishop 1:Kavu 2:Horror 3:Crab 4:Drake<br>
-- there are 2 copies of Bishop of Wings on the battlefield under the opponents control. These copies are hacked to have the creature types Kavu and Horror. They give 4 life whenever a Crab enters and create a Drake whenever a Crab dies.

1x Dralnu 1:Horror 2:Druid<br>
-- there is 1 copy of Dralnu's Crusade on the battlefield that causes Horrors to get +1/+1 and be Druids in addition to their other types. These modifications are applied in the order the Dralnus appear in the description. Layers and dependencies are not implemented.

1x Bishop Control<br>
-- This Bishop is on the battlefield under our control. This will be sacrificed to power the initial Soulblast.

1x Vanilla Swap<br>
-- This Golem is the target of the top Audacious Swap on the stack. It will be exiled to cast the initial Soulblast.

1x Vanilla Arcbond Blast<br>
-- This Golem has the Arcbond effect and will be the target of the initial Soulblast.

1x Bishop 4:Golem VIP<br>
-- This Bishop creates Golems and is marked as a VIP. That means the simulation is aborted if it dies. We need this bishop to profit from the computation!

1x Vanilla 1:Crab <Heartbeat<br>
-- This 3/3 Crab belongs to the display group "Heartbeat".

2x Bishop <Input >Output<br>
-- These bishops belong to the disply groub "Input". The spririts they create belong to the display group "Output"

1x Bishop Loud<br>
-- Whenever this bishop creates tokens during the simulation a summary of the boardstate will be printed.

# Boardstate summaries:

We print short information about the display groups defined in the input. We show the number of creatures in the group and the minimum life (toughness - damage) among creatures in the group.

# Boardstate requirements

To start a useful computation we need exactly one creature with the Swap keyword. Exactly one other creature needs to have the Blast keyword. That creature should also have the Arcbond keyword. More creatures can have Arcbond. We also need some creatures with Control to set the size of the Souldblast.

To profit from a computation we need 4 Vanillas to survive. Those are the targets of the remaining Audacious swaps on the stack. We also need a Bishop that creates Golems to survive. Mark these creatures with VIP.

The free Vanillas in the cost computation are the 5 Vanillas that are targeted by Audacious Swap copies, one of which gets exiled to cast Soulblast. The 6th free Vanilla is the target of the original Audacious Swap, that didn't get a copy tergeting it.

