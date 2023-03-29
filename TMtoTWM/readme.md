A short go program that generates [The Waterfall Model](http://nethack4.org/esolangs/waterfall/) programs for turing machines.

Input (via stdin):<br>
-- Line 1: A turing machine in [Standard TM Text Format](https://www.sligocki.com/2022/10/09/standard-tm-format.html) with n states and m symbols.<br>
-- Line 2: Initial input string for the TM

Output (via stdout):<br>
-- A Waterfall Model program with 8 + 4*m + n*m clocks, that runs the turing machine.

The first clock encodes the half-tape to the left of the TM head.<br>
The last clock encodes the half-tape to the right of the TM head.<br>
We interpret the half-tapes as a base m number with the least significant digit closest to the TM head. To write to a half-tape we multiply by m and add the new symbol. To read from a half-tape we divide by m and the remainder tells us which symbol we read.

The central clocks store the TM-transitions. When a transition clock empties it sets up the computation order to perform the instructions given by the corresponding TM-transition.