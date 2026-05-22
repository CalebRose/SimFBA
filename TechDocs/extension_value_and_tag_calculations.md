Calculating Extension Values
Note: when referring to current Values, reference the Value of the entire contract, not just the remaining years (should now be saved as a separate variable with the player)

1. Assign all players a position based on where they played the most snaps (above a certain floor), with primary position as the default
   a. Speed Rusher DEs and Pass Rush OLBs are grouped together in an “EDGE” category
   b. DTs and all other DEs are grouped together
   c. ILBs and all other OLBs are grouped together
2. For all players age 25 and younger, treat their Overall as being 4 points higher than it currently is
   The following steps should be applied to each position group separately
3. Within each position group, rank players from highest Overall to lowest, using Age as a secondary metric (ranked lowest to highest)
4. For players ranked at the top of their list, set their Value expectation equal to the Max Value of any current contract from that position group, PLUS 10%. If Overall is tied for the 5th place player, include all players with the same Overall in this category
   a. Exception: if the highest-Value contract is at least 150% of second-highest, exclude it from the list of current contracts to compare against (Caruso-Bordewyk Rule)
   b. The top tier for each position is the following
   i. 5: QB/RB/FB/TE/C/FS/SS/K/P
   ii. 10: OT/OG/DT/DE/OLB/ILB

5. For players not included in #4, create a custom set of current values to compare against. Use all players at the current position except:
   a. Do not consider players on rookie contracts
   b. Do not consider players on UDFA contracts (soon to be irrelevant)
   c. Do not consider contracts of players above a certain age
   i. QB/K: 34 and older
   ii. RB/FB: 27 and older
   iii. All other: 29 and older

The idea here is to prevent declining veterans at the end of long deals from inflating the expectations of younger mid-tier players

6. For all players in step #5, set their Value expectation equal to the highest current Value among players of equal or lower Overall in the custom set
7. Smoothing: create a second expectation Value for each player using the outputs from step 4-6. This second number is calculated as the average of the Values from the players who are [Overall +1] and [Overall -1] compared to the current player.

This took some manual fixing from me since there could be numbers missing, especially at the high end of each group

8. Set a player’s Value expectation as the higher number between steps 4-7
9. Adjust the Value based on the Age factor, using the table appropriate for that position (see appendix). This adjusted number is the final Value expectation for that player
10. Set AAV expectation as 40% of the Value expectation

Calculating Tag Values

1. Using the same position groups from step 1 above, rank all players within each group based on current year pay (NOT contract Value), calculated as current Bonus + current Salary
   a. If a player was traded, treat their Bonus as being what it was before they were traded (this also took a lot of manual work to restore Bonus amounts)
2. Calculate tag amounts based on the following averages of current year pay:
   a. Franchise: average of top 5
   b. Transition: average of top 10
   c. Playtime: average of 3rd – 20th
   d. Basic: average of 3rd – 25th

Note: the manual tags that teams can apply to any player during extension season (to automatically add one year to their current contract) always use the Franchise amount. The other tag amounts are only used to calculate 5th year option costs for 1st round draft picks. The appropriate tag is determined automatically using the following criteria:

1. Franchise: multiple Pro Bowl selections
2. Transition: one Pro Bowl selection
3. Playtime: starter-level snaps (765) in at least 2 of his first 3 seasons
4. Basic: none of the above

 
Appendix: Value Reduction by Age (and Position)

QB/K
Age Factor
23 15%
24 10%
25 5%
26 0%
27 -5%
28 -10%
29 -15%
30 -20%
31 -25%
32 -30%
33 -45%
34 -60%
35 -75%
36+ -90%

RB/FB
Age Factor
23 10%
24 5%
25 0%
26 -10%
27 -20%
28 -30%
29 -45%
30 -60%
31 -75%
32+ -90%

All Other
Age Factor
23 15%
24 10%
25 5%
26 0%
27 -10%
28 -20%
29 -30%
30 -40%
31 -50%
32 -60%
33 -75%
34+ -90%
