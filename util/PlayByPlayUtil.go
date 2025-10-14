package util

func GetOffensiveFormationByEnum(offForm uint8) string {
	offMap := getOffensiveFormationMap()
	return offMap[offForm]
}

func getOffensiveFormationMap() map[uint8]string {
	return map[uint8]string{
		0:  "N/A",
		1:  "Big Pistol",
		2:  "Big Spread Gun",
		3:  "Double Tight",
		4:  "Double Wing",
		5:  "Double Wing Strong",
		6:  "Double Wing Wide",
		7:  "Double Wing Spread",
		8:  "Double Wing Split",
		9:  "Empty Gun",
		10: "Flexbone",
		11: "Flexbone Strong",
		12: "Flexbone Wide",
		13: "Flexbone Gun",
		14: "Flexbone Gun Wide",
		15: "I Formation",
		16: "I Formation Heavy",
		17: "Near/Far",
		18: "Singleback",
		19: "Singleback Gun",
		20: "Splitback Gun",
		21: "Pistol",
		22: "Power Pistol",
		23: "Heavy Power Pistol",
		24: "Spread Pistol",
		25: "Spread",
		26: "Spread Gun",
		27: "Pony Spread Gun",
		28: "Wing-T",
		29: "Wing-T Split",
		30: "Wing-T Double Tight",
		31: "Wing Split Back Gun",
		32: "Wishbone",
		33: "Wishbone Strong",
		34: "Wishbone Wide",
	}
}

func GetDefensiveFormationByEnum(defForm uint8) string {
	defMap := getDefensiveFormationMap()
	return defMap[defForm]
}

func getDefensiveFormationMap() map[uint8]string {
	return map[uint8]string{
		0:  "N/A",
		1:  "3-2-6 Big Penny",
		2:  "3-2-6 Penny",
		3:  "3-2-6 Dime",
		4:  "3-3-5 Base",
		5:  "3-3-5 Nickel",
		6:  "3-3-5 Over",
		7:  "3-4 Bronco",
		8:  "3-4 Eagle",
		9:  "3-4 Okie",
		10: "3-4-4 Heavy",
		11: "4-1-6 Big Dime",
		12: "4-1-6 Dime",
		13: "4-2-5 Base",
		14: "4-2-5 Nickel",
		15: "4-2-5 Over",
		16: "4-3 Base",
		17: "4-3 Heavy",
		18: "4-3 Light",
		19: "4-3 Over",
		20: "4-4 Base",
		21: "4-4 Heavy",
		22: "4-4 Jumbo",
		23: "4-4 Over",
		24: "4-4 Under",
	}
}

func GetDefensiveTendencyByEnum(defTen uint8) string {
	defMap := getDefensiveTendencyMap()
	return defMap[defTen]
}

func getDefensiveTendencyMap() map[uint8]string {
	return map[uint8]string{
		0: "N/A",
		1: "Run Defense",
		2: "Pass Defense",
	}
}

func GetPlayTypeByEnum(playType uint8) string {
	playMap := getPlayTypeMap()
	return playMap[playType]
}

func getPlayTypeMap() map[uint8]string {
	return map[uint8]string{
		0: "Run",
		1: "Pass",
		2: "Punt",
		3: "Kickoff",
		4: "FG",
		5: "XP",
	}
}

func GetPlayNameByEnum(playName uint8) string {
	playNameMap := getPlayNameMap()
	return playNameMap[playName]
}

func getPlayNameMap() map[uint8]string {
	return map[uint8]string{
		0:  "N/A",
		1:  "Outside Run Left",
		2:  "Outside Run Right",
		3:  "Inside Run Left",
		4:  "Inside Run Right",
		5:  "Power Run Left",
		6:  "Power Run Right",
		7:  "Draw Left",
		8:  "Draw Right",
		9:  "Read Option Left",
		10: "Read Option Right",
		11: "Speed Option Left",
		12: "Speed Option Right",
		13: "Inverted Option Left",
		14: "Inverted Option Right",
		15: "Triple Option Left",
		16: "Triple Option Right",
		17: "Choice Inside Right",
		18: "Choice Inside Left",
		19: "Choice Outside Right",
		20: "Choice Outside Left",
		21: "Choice Power Right",
		22: "Choice Power Left",
		23: "Peek Inside Right",
		24: "Peek Inside Left",
		25: "Peek Outside Right",
		26: "Peek Outside Left",
		27: "Peek Power Right",
		28: "Peek Power Left",
		29: "Short",
		30: "Medium",
		31: "Long",
		32: "Screen",
		33: "Play Action Medium",
		34: "Play Action Long",
	}
}

func GetPointOfAttackByEnum(poa uint8) string {
	playMap := getPointOfAttackMap()
	return playMap[poa]
}

func getPointOfAttackMap() map[uint8]string {
	return map[uint8]string{
		0:  "N/A",
		1:  "Inside Left",
		2:  "Inside Right",
		3:  "Outside Left",
		4:  "Outside Right",
		5:  "Power Left",
		6:  "Power Right",
		7:  "Draw Left",
		8:  "Draw Right",
		9:  "Screen",
		10: "Quick",
		11: "Short",
		12: "Long",
	}
}

func GetPenaltyByEnum(penalty uint8) string {
	penaltyMap := getPenaltyMap()
	return penaltyMap[penalty]
}

func getPenaltyMap() map[uint8]string {
	return map[uint8]string{
		0:  "N/A",
		1:  "Catch Interference",
		2:  "Defensive Holding",
		3:  "Defensive Pass Interference",
		4:  "Delay of Game",
		5:  "Enchroachment",
		6:  "Face Mask",
		7:  "False Start",
		8:  "Holding",
		9:  "Holding Kicking Team",
		10: "Holding Returning Team",
		11: "Horse Collar",
		12: "Horse Collar Tackle",
		13: "Illegal Block",
		14: "Illegal Block Above the Waist",
		15: "Illegal Contact",
		16: "Illegal Double-Team Block",
		17: "Illegal Fair Catch Signal",
		18: "Illegal Formation",
		19: "Illegal Forward Pass",
		20: "Illegal Motion",
		21: "Illegal Shift",
		22: "Illegal Touch(Player Out of Bounds)",
		23: "Illegal Use of Hands",
		24: "Ineligible Downfield",
		25: "Intentional Grounding",
		26: "Kickoff Out of Bounds",
		27: "Neutral Zone Infraction",
		28: "Offensive Holding",
		29: "Offensive Pass Interference",
		30: "Offside",
		31: "Offsides",
		32: "Roughing the Kicker",
		33: "Roughing the Passer",
		34: "Running Into the Kicker",
		35: "Too Many Men on the Field",
		36: "Unnecessary Roughness",
		37: "Unsportsmanlike Conduct",
	}
}

func GetInjuryByEnum(injr uint8) string {
	injMap := getInjuryMap()
	return injMap[injr]
}

func getInjuryMap() map[uint8]string {
	return map[uint8]string{
		0:  "N/A",
		1:  "Achilles Tendonitis",
		2:  "ACL Bruise",
		3:  "ACL Tear",
		4:  "ACL Tendonitis",
		5:  "Ankle Bruise",
		6:  "Ankle Sprain",
		7:  "Back Disk Tear",
		8:  "Biceps Tear",
		9:  "Bruised Achilles",
		10: "Bruised Elbow",
		11: "Bruised Foot",
		12: "Bruised Hamstring",
		13: "Bruised Hip",
		14: "Bruised Thumb",
		15: "Bruised Toe",
		16: "Calf Tear",
		17: "Concussion",
		18: "Dislocated Ankle",
		19: "Dislocated Elbow",
		20: "Dislocated Foot",
		21: "Dislocated Shoulder",
		22: "Dislocated Thumb",
		23: "Dislocated Toe",
		24: "Elbow Tendonitis",
		25: "Fractured Ankle",
		26: "Fractured Foot",
		27: "Fractured Hip",
		28: "Fractured Jaw",
		29: "Fractured Ribs",
		30: "Fractured Spine",
		31: "Fractured Thumb",
		32: "Fractured Toe",
		33: "Fractured Wrist",
		34: "Groin Tear",
		35: "Hamstring Tendonitis",
		36: "High Ankle Sprain",
		37: "Hip Strain",
		38: "Hyperextended Back",
		39: "Illness",
		40: "Knee Meniscus Bruise",
		41: "Knee Meniscus Tear",
		42: "Lacerated Spleen",
		43: "MCL Bruise",
		44: "MCL Tear",
		45: "MCL Tendonitis",
		46: "Neck Bruise",
		47: "Patellar Tendon Bruise",
		48: "Patellar Tendon Tear",
		49: "Patellar Tendonitis",
		50: "PCL Bruise",
		51: "PCL Tear",
		52: "PCL Tendonitis",
		53: "Pectoral Tear",
		54: "Pulled Biceps",
		55: "Pulled Calf",
		56: "Pulled Groin",
		57: "Pulled Hamstring",
		58: "Pulled Pectoral",
		59: "Pulled Quadriceps",
		60: "Pulled Triceps",
		61: "Quadriceps Tear",
		62: "Rotator Cuff Tear",
		63: "Ruptured Achilles",
		64: "Ruptured Hamstring",
		65: "Separated Shoulder",
		66: "Shoulder Tendonitis",
		67: "Sprained Elbow",
		68: "Sprained Foot",
		69: "Sprained Knee",
		70: "Sprained Neck",
		71: "Sprained Rotator Cuff",
		72: "Sprained Thumb",
		73: "Sprained Toe",
		74: "Sprained Wrist",
		75: "Stinger",
		76: "Strained Back",
		77: "Strained Biceps",
		78: "Strained Calf",
		79: "Strained Groin",
		80: "Strained Hip",
		81: "Strained Pectoral",
		82: "Strained Quadriceps",
		83: "Strained Rotator Cuff",
		84: "Strained Shoulder",
		85: "Strained Triceps",
		86: "Triceps Tear",
		87: "Turf Toe",
		88: "Wrist Bruise",
	}
}

func GetInjuryLength(injr int) string {
	if injr == -4 {
		return "a quarter"
	}
	if injr == 2 {
		return "the remainder of the half"
	}
	if injr == 0 {
		return "N/A"
	}
	if injr == 1 {
		return "the rest of the game"
	}
	if injr == 2 {
		return "the rest of the Game and next week"
	}
	if injr > 2 && injr < 5 {
		return "a couple of weeks"
	}
	if injr < 8 {
		return "several Weeks"
	}
	return "most likely the remainder of the season"
}

func GetInjurySeverity(sev int) string {
	if sev == 0 {
		return "N/A"
	}
	if sev == 1 {
		return "Minor"
	}
	if sev == 2 {
		return "Moderate"
	}
	if sev == 3 {
		return "Severe"
	}
	if sev == 4 {
		return "Season=Ending"
	}
	return "Career-Ending"
}

func GetCoverageStr(cov bool) string {
	if cov {
		return "Zone"
	}
	return "Man"
}

func GetPuntVerb() string {
	list := []string{" punts for ", " punts it for ", " punts it off for ", " sends it away for "}
	return PickFromStringList(list)
}

func GetKickoffVerb(num int) string {
	var list []string
	if num == 1 {
		// First verb
		list = []string{" kicks for ", " kicks it off for "}
	} else if num == 2 {
		list = []string{" Fielded at the goalline by ", " there to catch it by the goalline is "}
	} else if num == 3 {
		list = []string{" Fielded at the goalline by ", " caught in the endzone by ", " what a kick! Fielded deep in the endzone by "}
	}

	return PickFromStringList(list)
}

func GetFairCatchStatement() string {
	list := []string{" And the receiver raises their hand, signaling a fair catch. ", " calls for a fair catch ", " waves for a fair catch. "}
	return PickFromStringList(list)
}

func GetFGStartingStatement(isFG bool) string {
	verb := "field goal"
	if !isFG {
		verb = "extra point"
	}
	list := []string{
		" with the " + verb + " attempt...",
		" with the " + verb + "...",
		verb + " kick is up...",
	}

	return PickFromStringList(list)
}

func GetBlockedStatement(isKick bool) string {
	verb := "kick"
	if !isKick {
		verb = "punt"
	}
	list := []string{
		" BLOCKED BY ",
		"...and it's BLOCKED! ",
		"...and the attempt is BLOCKED! Swatted away, and the defense scrambles for the ball! Block made by ",
		"...and the attempt is BLOCKED! Swatted down by ",
		"...and the attempt is BLOCKED! Denied by ",
		"...and it's BLOCKED! The crowd goes wild as the defense storms the field! ",
		" - and the ball is BLOCKED! What a turn of events! Blocked by ",
		"...and the attempt is BLOCKED! Swatted away, and the defense takes control! ",
		"...and the attempt is BLOCKED! The special teams unit has made a statement! ",
		"...and the attempt is BLOCKED! That's a game changer right there by ",
		"...thwarted! The " + verb + " is BLOCKED and the momentum shifts! ",
		"...and the attempt is SMOTHERED! The defense has put up the wall! Blocked by ",
		"...and the " + verb + " is REJECTED! A colossal block by ",
		"...denied! The " + verb + " is BLOCKED, sending shockwaves through the stadium! Block made by ",
		"...and it's STUFFED! The " + verb + " is BLOCKED and the defense is ecstatic! Denial made by ",
		"...and it's REBUFFED at the line! The " + verb + " is BLOCKED, incredible effort by ",
		"...and it's SNATCHED from the sky! BLOCKED with authority by ",
		"...and it's SLAPPED DOWN! The " + verb + " is BLOCKED - a crushing blow by ",
	}

	return PickFromStringList(list)
}

func GetFGEndStatement(isGood bool, isLeft, isOffUpright, isRight bool) string {
	list := []string{
		" and it's good!",
		" and it looks good!",
		" good.",
		" the kick is good.",
		" and it goes right through the uprights!",
		" the kick is good!",
	}
	if !isGood {
		list = []string{
			" and it's no good! ",
			" no good. ",
			" and the crowd groans, it's not through! No good. ",
			" and it's a miss! No good. ",
			" and the kick is wide! No good. ",
			" and the attempt falls short! No good. ",
		}
		if isLeft {
			list = append(list,
				" and its just slightly misses from the left. No good. ",
				" and it looks like it's just too far left. No good. ",
				" and it's too wide left. No good. ",
				" and it just misses from the left. The kicker looks frustrated. ",
				" and it veers left, no good! ",
				" and it hooks to the left! No good! ",
				" and the kick is left of the mark! No good. ",
				" and it's wide to the left, no good! ",
				" and it's pushed to the left! ",
				" and the ball misses to the left! ",
				" and it's off course, to the left! ",
			)
		} else if isOffUpright {
			list = append(list,
				" and it hits the upright, no good! ",
				" and it clangs off the upright! No good. ",
				" and the ball ricochets off the upright! No good. ",
				" and it's off the iron, denied! No good. ",
			)
		} else if isRight {
			list = append(list,
				" and it's wide to the right, no good! ",
				" and it's pushed to the right! ",
				" and the ball misses to the right! ",
				" and it's off course, to the right! ",
				" and its just slightly misses from the right. No good. ",
				" and it looks like it's just too far right. No good. ",
				" and it's too wide right. No good. ",
				" and it just misses from the left. The kicker looks frustrated. ",
				" and it veers right, no good! ",
				" and it hooks to the right! No good! ",
				" and the kick is right of the mark! No good. ",
			)
		}
	}

	return PickFromStringList(list)
}

func GetTouchbackStatement() string {
	list := []string{" Touchback. ", " And it looks like it will be a touchback, folks. ", " And it's a touchback, folks. "}
	return PickFromStringList(list)
}

func GetReturnVerb(yards int, touchdown, oob bool) string {
	list := []string{" returns the ball ", " runs it for ", " runs it by for ", " advances the ball for ", " hustles for "}
	if !touchdown && !oob {
		if yards > 9 {
			list = append(list,
				" spins past a defender, runs toward the left side and makes some headway! Bulldozes a defender but is quickly pulled down. Wow. Ran for ",
				" spins past a defender, runs toward the right side. Bulldozes a defender but is quickly pulled down. Wow. Ran for ",
				" charges upfield, shrugging off tacklers and marking a significant run of ",
				" blazes down the field, eluding defenders and finally tackled after a fantastic run of ",
				" blazes down the field, eluding defenders and finally tackled. Fantastic gain of ",
				" gallops downfield like a freight train, eventually brought down after a tremendous run of ",
				" gallops downfield like a freight train, eventually brought down after a tremendous effort. A run for ")
		} else if yards > 4 {
			list = append(list, " runs by the sideline, passes a defender... and now he's tackled. Ran for ",
				" darts down the middle of the field, jukes past a defender and is tackled. Goes for ",
				" finds a gap, dodges a tackle, and is eventually stopped. A run of ",
				" breaks to the outside, gaining momentum for a solid return of ",
				" hurdles over a defender and pushes forward, a decent return of ",
				" catches, tucks, and rolls out an agile run, netting ",
				" takes off like a cannonball, squeezing out a valuable ")
		} else if yards > 2 {
			list = append(list, " runs past a defender before getting tackled. Looks like it will be for ",
				" speeds past a defender but is swiftly tackled. Ran for ",
				" weaves through traffic is brought down after gaining ",
				"pushes ahead and is taken down after a small advance of ")
		} else {
			list = append(list, " makes a run for it but is brought down quickly. Ran for ",
				" speeds past a defender but is swiftly tackled. Ran for ",
				" scoots up the field for ",
				" squeezes past the line of scrimmage but is quickly halted. A gain of ",
				" charges forward but meets a wall of defenders. A short gain of ",
				" finds little room and is stopped after a modest gain of ",
				" stumbles slightly but regains footing, a short return for ",
				" gets a grip on the ball and scurries ahead, only managing ",
				" gathers the kick and charges ahead into the arms of the opposition, a minimal gain of ",
			)
		}
	} else if !touchdown && oob {
		if yards > 9 {
			list = append(list,
				" bolts to the edge and turns upfield with a burst of speed, pushing his way out of bounds after gaining ",
				" dashes toward the sideline, narrowly avoiding a tackle and stepping out after a solid ",
				" evades the first wave of tacklers and smartly steps out to avoid a hit. A good return for ",
			)
		} else if yards > 4 {
			list = append(list,
				" shuffles his way to the side, squeezing out a few hard-earned yards before being forced out at ",
				" takes it to the sideline, sidestepping a defender and going out with a gain of ",
				" manages to skirt the sideline, picking up a modest ",
			)
		} else {
			list = append(list, " nimbly steps along the sideline before being pushed out of bounds after ",
				" speeds along the edge, narrowly avoiding tackles and stepping out after ",
				" tiptoes the sideline, making good ground before being ushered out after ",
				" is quickly corralled by the coverage team and nudged out at the ",
				" sprints to the boundary, only to be met by a defender. Ends up with just ",
				" catches the ball and instinctively heads out of bounds. Not much there, just ",
			)
		}
	} else {
		list = append(list,
			" darts toward the sideline, turns upfield and takes it all the way! Touchdown! ",
			" races across the field, leaving defenders in his wake and crosses the goal line for a touchdown! ",
			" explodes past the coverage team, finds an opening and sprints to the end zone for six points! ",
			" finds a seam and explodes through, he could go all the way, touchdown after a tremendous ",
			" catches the kick, weaves through traffic like a speedboat through water, and he's gone! ",
			" grabs the ball and turns on the jets, leaving defenders in the dust. That's six points with an electric ",
		)
	}

	return PickFromStringList(list)
}

func GetScrambleText(scramble bool) string {
	if scramble {
		return " scrambles, and"
	}
	return ""
}
