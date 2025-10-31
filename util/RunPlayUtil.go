package util

func GetRunVerb(yards int, playName, poa string, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) string {
	list := []string{}
	if !touchdown && !outOfBounds && !twoPtConversion && !fumble && !safety {
		list = append(list, " carries for ", " runs the ball ", " runs it for ", " runs it by for ", " advances the ball for ", " hustles for ")
	}
	// Basic run
	switch playName {
	case "Inside Run Left", "Inside Run Right":
		list = append(list, getInsideRunList(yards, playName == "Inside Run Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Outside Run Left", "Outside Run Right":
		list = append(list, getOutsideRunList(yards, playName == "Outside Run Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Power Run Left", "Power Run Right":
		list = append(list, getPowerRunList(yards, playName == "Power Run Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Draw Left", "Draw Right":
		list = append(list, getDrawRunList(yards, playName == "Draw Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Read Option Left", "Read Option Right":
		list = append(list, getReadOptionList(yards, playName == "Read Option Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Speed Option Left", "Speed Option Right":
		list = append(list, getSpeedOptionList(yards, playName == "Speed Option Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Inverted Option Left", "Inverted Option Right":
		list = append(list, getInvertedOptionList(yards, playName == "Inverted Option Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Triple Option Left", "Triple Option Right":
		list = append(list, getTripleOptionList(yards, playName == "Triple Option Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Choice Inside Left", "Choice Inside Right":
		list = append(list, getInsideRunList(yards, poa == "Inside Run Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Choice Outside Left", "Choice Outside Right":
		list = append(list, getOutsideRunList(yards, poa == "Outside Run Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Peek Inside Left", "Peek Inside Right":
		list = append(list, getInsideRunList(yards, poa == "Inside Run Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Peek Outside Left", "Peek Outside Right":
		list = append(list, getOutsideRunList(yards, poa == "Outside Run Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	case "Peek Power Left", "Peek Power Right":
		list = append(list, getPowerRunList(yards, poa == "Power Run Left", touchdown, outOfBounds, twoPtConversion, fumble, safety)...)
	}

	return PickFromStringList(list)
}

func getPowerRunList(yards int, isleft, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) []string {
	list := []string{}
	direction := "left"
	if isleft {
		direction = "right"
	}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	// Depending on the yards gained or lost, append appropriate descriptions
	if touchdown {
		list = append(list, " bulldozes through the line and powers his way into the end zone for a TOUCHDOWN! ",
			" charges forward, pushing defenders back and breaking the plane for a score! TOUCHDOWN! ",
			" showcases sheer strength, fighting off tacklers and plunging into the end zone! TOUCHDOWN! ",
			" demonstrates a powerful run, shedding tackles and reaching paydirt! TOUCHDOWN! ",
			" steamrolls his way through the defense and into the end zone for six! TOUCHDOWN! ",
			" plows over the goal line, dragging defenders with him for the TOUCHDOWN! ")
	} else if outOfBounds {
		list = append(list, " makes a dash to the "+direction+" and steps out after"+gainStatement,
			" stretches the play to the "+direction+" before being forced out,"+gainStatement,
			" barrels to the "+direction+" and is pushed out after"+gainStatement,
			" powers his way to the sideline, gaining valuable yards before stepping out ",
			" shows off his strength, pushing through to the "+direction+" before going out at ",
			" maintains a strong run to the "+direction+" and exits the field after"+gainStatement)
	} else if twoPtConversion && touchdown {
		list = append(list, " pushes through the line for the crucial two-point conversion! ",
			" barrels into the end zone for the two-point conversion! ",
			" dives across the goal line, and the two-point attempt is good! ",
			" finds a tiny crease and converts for two points! ",
			" outmuscles the defenders and powers in for the conversion! ",
			" drives through the defenders and secures the two-point conversion with power! ",
			" leans into the defense, pushing forward for the successful conversion! ",
			" fights through the line with a surge of strength for two points! ",
			" delivers a forceful run, crossing into the end zone for the conversion! ")
	} else if twoPtConversion && !touchdown {
		list = append(list, " can't push through the line for the crucial two-point conversion! No good! ",
			" barrels into a lineman and misses the end zone for the two-point conversion! No good! ",
			" makes a dive, but misses the goal line. The two-point attempt is no good! ",
			" can't find a gap and is brought down. Two point conversion is no good! ",
			" is outmuscled by the defenders. The two point conversion is no good! ",
			" attempts to bulldoze through but is halted before the line. No conversion! ",
			" tries to overpower the defense but falls short of the end zone! ",
			" pushes with all his might but can't break the plane. Conversion fails! ",
			" meets a wall of defenders and can't punch it through for the two points! ")
	} else if fumble {
		list = append(list, " fights for extra yards but loses the ball in the skirmish after"+gainStatement,
			" powers forward but loses control of the football, fumbling after"+gainStatement,
			" has the ball jarred loose after"+gainStatement,
			" makes a costly mistake, fumbling the ball away after ",
			" can't hold on to the pigskin, and it's on the ground after ")
	} else if safety {
		list = append(list, " is caught in the end zone for a safety! ",
			" is wrapped up in the end zone, and it's two points the other way! ",
			" gets trapped in the end zone during a power run, resulting in a safety! ",
			" can't break free from the defense in the end zone, conceding a safety! ",
			" is swarmed by the defense for a safety! ",
			" tries to get out of the end zone but is trapped for a safety! ",
			" can't escape the clutches of the defense – that's a safety! ")
	} else {
		switch {
		case yards > 14:
			list = append(list, " bursts through the "+direction+" and breaks free for"+gainStatement,
				" passes the line of scrimmage, dodges right past a linebacker and makes a break for it before being "+tackleVerb+" for ",
				" charges upfield, shrugging off tacklers and marking"+gainStatement,
				" blazes down the field, eluding defenders and finally "+tackleVerb+" after"+gainStatement,
				" blazes down the field, eluding defenders and finally "+tackleVerb+". Fantastic gain of ",
				" chugs downfield like a freight train, eventually "+tackleVerb+" after"+gainStatement,
				" gallops downfield like a gazelle, eventually "+tackleVerb+" after a tremendous effort. A run for ",
				" steamrolls through defenders, showcasing brute strength for"+gainStatement,
				" bulldozes his way downfield, an unstoppable force for"+gainStatement,
				" plows through the line like a hot knife through butter,"+gainStatement,
				" demonstrates sheer power, breaking multiple tackles on his way to"+gainStatement,
				" turns a routine power run into a highlight reel, thundering down the field for"+gainStatement,
				" embodies the power back archetype, trampling over defenders for"+gainStatement)
		case yards > 9:
			list = append(list,
				" powers through the "+direction+" for a strong run,"+gainStatement,
				" finds a hole in the "+direction+" and hits it with authority,"+gainStatement,
				" darts through the "+direction+" side, with a burst of speed,"+gainStatement,
				" carves up the defense on the "+direction+" for"+gainStatement,
				" exploits a gap to the "+direction+" and rumbles past defenders for"+gainStatement,
				" turns on the afterburners and leaves the defense in the dust,"+gainStatement,
				" weaves through the defense to the "+direction+" and breaks away,"+gainStatement,
				" finds daylight on the "+direction+", turning on the jets for"+gainStatement,
				" blasts past the line into the "+direction+" side, barreling forward for"+gainStatement,
				" rumbles ahead with determination, pushing defenders aside for"+gainStatement,
				" showcases a blend of power and agility, turning a tough run into ",
				" powers ahead, dragging defenders with him for"+gainStatement,
				" barrels down the "+direction+", shrugging off tackles for"+gainStatement,
				" exhibits a classic power run, plowing forward to pick up ",
				" charges ahead with sheer force, carving out a nice chunk of yards for"+gainStatement)
		case yards > 4:
			list = append(list,
				" muscles his way past the line and into the "+direction+","+gainStatement,
				" surges forward with power, shouldering through to the "+direction+" for ",
				" follows his blockers to the "+direction+" and powers through for"+gainStatement,
				" with a hop and a skip, sidesteps a defender and sprints to the "+direction+" for ",
				" hits the hole on the "+direction+" with gusto,"+gainStatement,
				" grinds out a hard-earned few yards, battling his way through for"+gainStatement,
				" fights for every inch, muscling ahead for ",
				" pushes the pile, demonstrating strength as he gains ",
				" gains some tough yards on the power run for"+gainStatement,
				" moves the chains bit by bit, a testament to his strength,"+gainStatement,
				" battles through the line, securing a few valuable yards for ")
		case yards > 0:
			list = append(list,
				" nudges forward on the "+direction+","+gainStatement,
				" makes a run for it but is brought down quickly. Ran for"+gainStatement,
				" speeds past a defender but is swiftly tackled. Ran for ",
				" scoots up the field for"+gainStatement,
				" squeezes past the line of scrimmage but is quickly halted. A gain of ",
				" charges forward but meets a wall of defenders. A short gain of ",
				" meets a wall of defenders,"+gainStatement,
				" barely nudges forward, a testament to the defense's strength,"+gainStatement,
				" finds little room and is stopped after"+gainStatement,
				" stumbles slightly but regains footing,"+gainStatement,
				" gets a grip on the ball and scurries ahead,"+gainStatement,
				" gathers the kick and charges ahead into the arms of the opposition,"+gainStatement)
		case yards == 0:
			list = append(list, " is met at the line of scrimmage, no gain on the play. ",
				" is immediately stopped at the line of scrimmage, no gain on the play. ",
				" is stonewalled at the line, no gain on the play. ",
				" finds no room to maneuver, stopped dead in his tracks. ",
				" is immediately tackled by the defensive line. No gain on the play. ")
		default: // Negative yardage
			list = append(list, " is wrapped up behind the line,"+gainStatement,
				" is stopped in the backfield for"+gainStatement,
				" meets a defender and is "+tackleVerb+" in the backfield,"+gainStatement,
				" tries to find a gap, and just like that he's "+tackleVerb+","+gainStatement,
				" is swallowed by the defense,"+gainStatement,
				" can't find a crease and is "+tackleVerb+" behind the line"+gainStatement,
				" attempts to find a gap and is "+tackleVerb+","+gainStatement)
		}
	}
	return list
}

func getOutsideRunList(yards int, isleft, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) []string {
	list := []string{}
	direction := "left"
	if isleft {
		direction = "right"
	}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	// Depending on the yards gained or lost, append appropriate descriptions
	if touchdown {
		list = append(list, " runs outside the line and bursts into the end zone for a TOUCHDOWN! ",
			" dashes around the defense and into the promised land for six! TOUCHDOWN! Ran ",
			" avoids the traffic and dances into the end zone! TOUCHDOWN! Ran ",
			" counters around line and sprints into the end zone untouched! TOUCHDOWN! Ran ",
			" turns on the speed and leaves everyone in his dust – that's a TOUCHDOWN! Ran ",
			" breaks an outside tackle, then another, and he's in for the score! TOUCHDOWN! Ran ")
	} else if outOfBounds {
		list = append(list, " makes a dash to the "+direction+" and steps out after"+gainStatement,
			" skirts the sideline before stepping out,"+gainStatement,
			" scampers to the edge and just tiptoes out of bounds,"+gainStatement,
			" races towards the sideline, securing a few extra yards before going out at ",
			" stretches the play to the "+direction+" before being forced out,"+gainStatement,
			" finds a lane to the "+direction+" but is pushed out after"+gainStatement,
			" accelerates to the "+direction+" and is nudged out,"+gainStatement,
			" cuts to the "+direction+", evading a tackle and stepping out with "+gainStatement,
			" pushes his way to the "+direction+" sideline and steps out,"+gainStatement,
			" showcases agility with a dash to the "+direction+", stepping out after"+gainStatement)
	} else if twoPtConversion && touchdown {
		list = append(list, " pushes around the line for the crucial two-point conversion! ",
			" barrels into the end zone for the two-point conversion! ",
			" makes a run for the goal line, and the two-point attempt is good! ",
			" avoids the dogpile and converts for two points! ",
			" outruns the defenders on the outside and powers in for the conversion! ",
			" finds a crease to the outside and crosses into the end zone for the two-pointer! ",
			" sweeps to the "+direction+" and breaks the plane for a crucial two points! ",
			" darts to the "+direction+" edge and secures the two-point conversion with ease! ",
			" capitalizes on the outside track and dives in for the additional two points! ")
	} else if twoPtConversion && !touchdown {
		list = append(list, " can't find a gap on the outside for the crucial two-point conversion! No good! ",
			" barrels into a lineman and misses the end zone for the two-point conversion! No good! ",
			" makes an outside dive, but misses the goal line. The two-point attempt is no good! ",
			" is stopped outside the line and is brought down. Two point conversion is no good! ",
			" is bottlenecked on the outside route, failing to convert the two points! ",
			" tries an outside leap but falls short of the end zone, two points denied! ",
			" is cornered to the "+direction+" and stopped, missing the critical two points! ",
			" aims for the corner but is tackled just shy of the goal line, no conversion! ")
	} else if fumble {
		list = append(list, " coughs up the ball after"+gainStatement,
			" loses the grip on the ball mid-run after"+gainStatement,
			" has the ball jarred loose after"+gainStatement,
			" makes a costly mistake, fumbling the ball away after"+gainStatement,
			" can't hold on to the pigskin, and it's on the ground after"+gainStatement,
			" navigates to the "+direction+" but loses control of the ball, dropping it after"+gainStatement,
			" is on a promising run to the "+direction+" but fumbles, the ball slips away after"+gainStatement,
			" makes a break to the "+direction+" but the ball is knocked loose, fumbling after"+gainStatement)
	} else if safety {
		list = append(list, " is caught in the end zone for a safety! ",
			" is wrapped up in the end zone, and it's two points the other way! ",
			" is swarmed by the defense for a safety! ",
			" tries to get out of the end zone but is trapped for a safety! ",
			" can't escape the clutches of the defense – that's a safety! ",
			" is ensnared in the end zone while trying an outside run, resulting in a safety! ",
			" attempts to escape the end zone via the "+direction+" but is tackled for a safety! ",
			" gets cornered in the end zone on an outside attempt, conceding two points! ",
			" retreats to the "+direction+" in the end zone but is caught, a costly safety! ")
	} else {
		switch {
		case yards > 14:
			list = append(list, " bursts through the "+direction+" and breaks free for a massive gain of ",
				" spins past a defender, runs toward the "+direction+" side and makes some headway! Bulldozes a defender but is quickly "+tackleVerb+". Wow. Ran for ",
				" spins past a defender, runs toward the "+direction+" side. Bulldozes a defender but is quickly "+tackleVerb+". Wow. Ran for ",
				" charges upfield, shrugging off tacklers and marking"+gainStatement,
				" blazes down the field, eluding defenders and finally "+tackleVerb+" after"+gainStatement,
				" blazes down the field, eluding defenders and finally "+tackleVerb+". Fantastic gain of ",
				" gallops downfield like a freight train, eventually "+tackleVerb+" after"+gainStatement,
				" turns on the gas and speeds off, eventually "+tackleVerb+" after a tremendous effort. A run for ")
		case yards > 9:
			list = append(list,
				" darts through the "+direction+" side, with a burst of speed,"+gainStatement,
				" carves up the defense on the outside "+direction+" for"+gainStatement,
				" runs outside to the "+direction+" and rumbles past defenders for "+gainStatement,
				" turns on the afterburners and leaves the defense in the dust,"+gainStatement,
				" weaves through the defense to the "+direction+" and breaks away,"+gainStatement,
				" finds daylight on the "+direction+", turning on the jets"+gainStatement,
				" runs around the line into the "+direction+" side, barreling forward for ")
		case yards > 4:
			list = append(list,
				" muscles his way around the line and into the "+direction+", securing ",
				" takes a handoff to the outside "+direction+" and barrels ahead for "+gainStatement,
				" follows his blockers to the outside "+direction+" and speeds through for "+gainStatement,
				" counters to the "+direction+" with gusto,"+gainStatement,
				" accelerates around the edge, turning upfield for"+gainStatement,
				" rounds the "+direction+" corner with a burst,"+gainStatement,
				" glides to the "+direction+", using his speed to notch ",
				" skirts the "+direction+" sideline, eluding tackles for"+gainStatement,
				" sweeps to the "+direction+", cutting upfield for"+gainStatement,
				" darts around the "+direction+" side, scampering for"+gainStatement,
				" showcases his agility, turning the corner on the "+direction+" for"+gainStatement,
			)
		case yards > 0:
			list = append(list,
				" nudges forward on the "+direction+","+gainStatement,
				" edges his way to the "+direction+" for"+gainStatement,
				" maneuvers around the "+direction+" end,"+gainStatement,
				" pushes around the "+direction+","+gainStatement,
				" stretches the play to the "+direction+","+gainStatement,
				" battles past the "+direction+" edge for a modest ",
				" skirts the fringe, squeezing out a few yards to the "+direction+" for "+gainStatement,
				" leans forward, fighting for every inch to the "+direction+","+gainStatement,
				" makes a run for it but is brought down quickly. Ran for ",
				" speeds past a defender but is swiftly tackled. Ran for ",
				" finds little room on the outside and is "+tackleVerb+" after"+gainStatement,
				" stumbles slightly but regains footing,"+gainStatement,
				" gets a grip on the ball and scurries ahead,"+gainStatement,
			)
		case yards == 0:
			list = append(list, " is met on the outside, no gain on the play. ",
				" is immediately "+tackleVerb+" on the outside, no gain on the play. ",
				" can't edge past the defensive line. No gain on the play. ",
				" tries to turn the corner on the "+direction+" but is stopped dead in his tracks, no gain.",
				" is corralled at the "+direction+" boundary, unable to advance, no gain.",
				" attempts an outside run to the "+direction+" but is swarmed, no progress made. ",
				" seeks an opening on the "+direction+" but is stymied at the line, no gain. ")
		default: // Negative yardage
			list = append(list, " is wrapped up behind the line, a loss of ",
				" is "+tackleVerb+" in the backfield for "+gainStatement,
				" meets a defender and is "+tackleVerb+" in the backfield for"+gainStatement,
				" is "+tackleVerb+" outside the line and brought backwards for"+gainStatement,
				" attempts to find a gap and is "+tackleVerb+" for"+gainStatement,
				" is cornered and "+tackleVerb+" behind the line, losing ",
				" attempts to swing to the "+direction+" but is dropped for a loss of ",
				" is ensnared by the defense on the "+direction+", a setback of ",
				" tries to make headway to the "+direction+" but retreats, losing ",
				" is overwhelmed in the backfield trying to go "+direction+", a loss of ",
				" can't turn upfield on the "+direction+" and is tackled for a loss of ")
		}
	}
	return list
}

func getInsideRunList(yards int, isleft, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) []string {
	list := []string{}
	direction := "left"
	if isleft {
		direction = "right"
	}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	// Depending on the yards gained or lost, append appropriate descriptions
	if touchdown {
		list = append(list, " finds a seam and bursts into the end zone for a TOUCHDOWN! ",
			" slices through the defense and into the promised land for six! TOUCHDOWN! Ran ",
			" zigzags his way through the traffic and dances into the end zone! TOUCHDOWN! Ran ",
			" blasts through the line and sprints into the end zone untouched! TOUCHDOWN! Ran ",
			" turns on the speed and leaves everyone in his dust – that's a TOUCHDOWN! Ran ",
			" breaks one tackle, then another, and he's in for the score! TOUCHDOWN! Ran ")
	} else if outOfBounds {
		list = append(list, " makes a dash to the "+direction+" and steps out after"+gainStatement,
			" skirts the sideline before stepping out,"+gainStatement,
			" scampers to the edge and just tiptoes out of bounds,"+gainStatement,
			" races towards the sideline, securing a few extra yards before going out at ",
			" stretches the play to the "+direction+" before being forced out,"+gainStatement)
	} else if twoPtConversion && touchdown {
		list = append(list, " pushes through the line for the crucial two-point conversion! ", " barrels into the end zone for the two-point conversion! ",
			" dives across the goal line, and the two-point attempt is good! ",
			" finds a tiny crease and converts for two points! ",
			" outmuscles the defenders and powers in for the conversion! ")
	} else if twoPtConversion && !touchdown {
		list = append(list, " can't push through the line for the crucial two-point conversion! No good! ",
			" barrels into a lineman and misses the end zone for the two-point conversion! No good! ",
			" makes a dive, but misses the goal line. The two-point attempt is no good! ",
			" can't find a gap and is brought down. Two point conversion is no good! ",
			" is outmuscled by the defenders. The two point conversion is no good! ")
	} else if fumble {
		list = append(list, " coughs up the ball after"+gainStatement,
			" loses the grip on the ball mid-run after"+gainStatement,
			" has the ball jarred loose after"+gainStatement,
			" makes a costly mistake, fumbling the ball away after ",
			" can't hold on to the pigskin, and it's on the ground after ")
	} else if safety {
		list = append(list, " is caught in the end zone for a safety! ",
			" is wrapped up in the end zone, and it's two points the other way! ",
			" is swarmed by the defense for a safety! ",
			" tries to get out of the end zone but is trapped for a safety! ",
			" can't escape the clutches of the defense – that's a safety! ")
	} else {
		switch {
		case yards > 14:
			list = append(list, " bursts through the "+direction+" and breaks free,"+gainStatement,
				" passes the line of scrimmage, dodges right past a linebacker and makes a break for it before being "+tackleVerb+" for "+gainStatement,
				" spins past a defender, runs toward the "+direction+" side and makes some headway! Bulldozes a defender but is quickly pulled down. Wow. Ran for "+gainStatement,
				" spins past a defender, runs toward the "+direction+" side. Bulldozes a defender but is quickly pulled down. Wow. Ran for "+gainStatement,
				" charges upfield, shrugging off tacklers and marking a significant run of ",
				" blazes down the field, eluding defenders and finally "+tackleVerb+" after"+gainStatement,
				" blazes down the field, eluding defenders and finally "+tackleVerb+". Fantastic gain of ",
				" gallops downfield like a freight train, eventually "+tackleVerb+" after"+gainStatement,
				" gallops downfield like a freight train, eventually "+tackleVerb+" after a tremendous effort,"+gainStatement)
		case yards > 9:
			list = append(list,
				" powers through the "+direction+" for a strong run,"+gainStatement,
				" finds a hole in the "+direction+" and hits it with authority,"+gainStatement,
				" darts through the "+direction+" side, with a burst of speed,"+gainStatement,
				" carves up the defense on the "+direction+" for"+gainStatement,
				" exploits a gap to the "+direction+" and rumbles past defenders for "+gainStatement,
				" turns on the afterburners and leaves the defense in the dust,"+gainStatement,
				" weaves through the defense to the "+direction+" and breaks away,"+gainStatement,
				" finds daylight on the "+direction+", turning on the jets for"+gainStatement,
				" blasts past the line into the "+direction+" side, barreling forward for "+gainStatement)
		case yards > 4:
			list = append(list,
				" squeezes through a gap on the "+direction+" and grinds out ",
				" muscles his way past the line and into the "+direction+","+gainStatement,
				" takes a handoff to the "+direction+" and barrels ahead for"+gainStatement,
				" follows his blockers to the "+direction+" and powers through for"+gainStatement,
				" with a hop and a skip, sidesteps a defender and sprints to the "+direction+" for ",
				" hits the hole on the "+direction+" with gusto,"+gainStatement,
				" jukes to the "+direction+" and finds a seam,"+gainStatement,
				" with quick feet, darts to the "+direction+" and snags ")
		case yards > 0:
			list = append(list,
				" nudges forward on the "+direction+", managing to get ",
				" makes a run for it but is brought down quickly. Ran for "+gainStatement,
				" speeds past a defender but is swiftly tackled. Ran for "+gainStatement,
				" scoots up the field for ",
				" squeezes past the line of scrimmage but is quickly halted. A gain of ",
				" charges forward but meets a wall of defenders. A short gain of ",
				" finds little room and is stopped after"+gainStatement,
				" stumbles slightly but regains footing,"+gainStatement,
				" gets a grip on the ball and scurries ahead,"+gainStatement,
				" gathers the kick and charges ahead into the arms of the opposition,"+gainStatement)
		case yards == 0:
			list = append(list, " is met at the line of scrimmage, no gain on the play. ",
				" is immediately stopped at the line of scrimmage, no gain on the play. ",
				" immediately bumps into a lineman. No gain on the play. ",
				" is immediately tackled by the defensive line. No gain on the play. ")
		default: // Negative yardage
			list = append(list, " is wrapped up behind the line,"+gainStatement,
				" is stopped in the backfield for"+gainStatement,
				" meets a defender and is tackled in the backfield for"+gainStatement,
				" tries to find a gap, and just like that he's stopped for"+gainStatement,
				" finds no gap, tackled for"+gainStatement,
				" finds no gap, brought down for"+gainStatement,
				" attempts to find a gap and is brought down for"+gainStatement)
		}
	}
	return list
}

func getDrawRunList(yards int, isleft, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) []string {
	list := []string{}
	direction := "left"
	if isleft {
		direction = "right"
	}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	// Depending on the yards gained or lost, append appropriate descriptions
	if touchdown {
		list = append(list, " deceives the defense with a draw and sprints into the end zone for a TOUCHDOWN! ",
			" executes the draw play to perfection, finding a path to the end zone for six! TOUCHDOWN! ",
			" capitalizes on the draw, darting through the confused defenders for a TOUCHDOWN! ")
	} else if outOfBounds {
		list = append(list, " takes the draw to the "+direction+" and steps out after gaining ",
			" takes the draw to the "+direction+" and steps out,"+gainStatement,
			" takes the draw to the "+direction+" and steps out out bounds,"+gainStatement,
			" takes the draw to the "+direction+" and steps to the sideline,"+gainStatement,
			" uses the draw to find space and then smartly steps out of bounds"+gainStatement)
	} else if twoPtConversion && touchdown {
		list = append(list, " fools the defense with a draw and secures the two-point conversion! ",
			" slips through the line on a draw play, converting for two points! ")
	} else if twoPtConversion && !touchdown {
		list = append(list, " attempts a draw for the two-point conversion but is stopped short! ",
			" can't find the end zone on the draw play, missing the conversion! ")
	} else if fumble {
		list = append(list, " fumbles the ball on the draw play after "+gainStatement,
			" loses the ball on the draw, turning it over after a brief gain of ")
	} else if safety {
		list = append(list, " is trapped in the end zone on a draw play, resulting in a safety! ",
			" can't escape the clutches of the defense on the draw, conceding a safety! ")
	} else {
		switch {
		case yards > 14:
			list = append(list, " fools everyone with a draw and breaks free for"+gainStatement,
				" fools the defense with a draw, breaking free for"+gainStatement,
				" turns the draw play into a highlight, dashing past the defense for"+gainStatement,
				" turns the draw play into a highlight, sprinting past the defense for"+gainStatement,
				" turns the draw play into a highlight, sprinting past"+gainStatement,
				" turns the draw play into a highlight, dashing past for"+gainStatement,
				" turns the draw play into a highlight, running past the defense for"+gainStatement,
				" turns the draw play into a highlight, running "+gainStatement)
		case yards > 9:
			list = append(list, " navigates through the line on a draw,"+gainStatement,
				" finds a seam on the draw and accelerates"+gainStatement,
				" finds a seam on the draw and runs"+gainStatement,
				" finds a seam on the draw and makes a dash for it,"+gainStatement)
		case yards > 4:
			list = append(list, " grinds out yards with a well-executed draw,"+gainStatement,
				" picks his way through the defense on the draw"+gainStatement,
				" finds a gap on the draw and makes a run for"+gainStatement,
				" finds a gap on the draw and makes"+gainStatement,
				" finds a gap on the draw,"+gainStatement)
		case yards > 0:
			list = append(list, " gains a few tough yards on the draw,"+gainStatement,
				" makes a small but positive gain on the draw,"+gainStatement,
				" makes a run for it on the draw,"+gainStatement,
				" attempts the draw and is "+tackleVerb+""+gainStatement)
		default: // Negative yardage
			list = append(list, " is swallowed up behind the line on the draw,"+gainStatement,
				" can't fool the defense on the draw,"+gainStatement,
				" can't find a gap on the draw and is "+tackleVerb+""+gainStatement)
		}
	}
	return list
}

func getReadOptionList(yards int, isleft, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) []string {
	list := []string{}
	direction := "left"
	if isleft {
		direction = "right"
	}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	// Depending on the yards gained or lost, append appropriate descriptions
	if touchdown {
		list = append(list, " makes the perfect read and takes it all the way to the house for a TOUCHDOWN! ",
			" reads the goalline defense, opts to keep it, and dashes to the end zone for six! TOUCHDOWN! ",
			" deceives the goalline defense with a slick read option, sprinting into the end zone for a score! TOUCHDOWN! ",
			" navigates the read option with precision, blazing his way to the end zone for a spectacular TOUCHDOWN! ",
			" outwits the goalline defense with a masterful read, darting to paydirt for the TOUCHDOWN! ",
			" showcases his athleticism on the goalline, turning the read option into a thrilling TOUCHDOWN run! ")
	} else if outOfBounds {
		list = append(list, " opts for the "+direction+" and steps out after a decent gain. ",
			" reads, keeps, and sprints to the "+direction+", going out of bounds after gaining yards. ",
			" leverages the read option to gain the edge and smartly steps out,"+gainStatement,
			" exploits the defense's hesitation on the read and scoots out of bounds for"+gainStatement,
			" exploits the defense's hesitation on the option and scoots out of bounds for"+gainStatement,
			" surprises the defensive line on the option and edges out of bounds for"+gainStatement)
	} else if twoPtConversion && touchdown {
		list = append(list, " executes the read option flawlessly on the goalline for a two-point conversion! ",
			" makes the right read in the red zone and converts for two crucial points! ",
			" brilliantly navigates the read option, crossing into the end zone for the extra points! ",
			" capitalizes on the read option, weaving his way for a successful two-point play! ")
	} else if twoPtConversion && !touchdown {
		list = append(list, " tries the read option for the conversion but gets stopped short! ",
			" opts to keep but can't break through for the two points. ",
			" hesitates on the read option, can't break through for the two points. ",
			" gives the read option a go for the two points, but the defense stands tall, denying the conversion. ",
			" opts for the read on the conversion attempt, yet the defense snuffs it out, thwarting the extra points. ",
			" stalls on the read option attempt and stopped short of the goalline, thwarting the extra points. ",
			" stalls on the read option attempt and stopped short of the end zone, thwarting the extra points. ",
			" stalls on the read option attempt and stopped short of the end zone, thwarting the two point attempt. ",
			" trips up on the read option attempt and stopped short of the goalline, thwarting the extra points. ",
			" trips up on the read option attempt and stopped short of the end zone, thwarting the extra points. ",
			" trips up on the read option attempt and stopped short of the end zone, thwarting the two point attempt. ",
			" attempts the read option and is stopped short of the goalline, thwarting the extra points. ",
			" attempts the read option and is stopped short of the end zone, thwarting the extra points. ",
			" attempts the read option and is stopped short of the end zone, thwarting the two point attempt. ")
	} else if fumble {
		list = append(list, " fumbles while executing the read option, losing the ball after gaining ",
			" makes a read but loses the ball on the run, a costly turnover after ",
			" tries to navigate the read option but ends up turning the ball over with a fumble after ",
			" keeps on the read option, only to lose control of the ball, squandering ")
	} else if safety {
		list = append(list, " gets trapped in the end zone on a read option, resulting in a safety! ",
			" opts to keep but is caught in the end zone for a safety! ",
			" disastrously misreads the option in the end zone, leading to a safety against his team. ",
			" hesitates on the option in the end zone, leading to a safety against his team. ",
			" hesitates on the read option in the end zone, leading to a safety against his team. ",
			" is enveloped by the defense in the end zone during the read option, resulting in a safety. ")
	} else {
		switch {
		case yards > 14:
			list = append(list, " makes a great read and breaks free,"+gainStatement,
				" reads the defense perfectly and darts through the "+direction+" for"+gainStatement,
				" takes the read option and finds a gap in the defense,"+gainStatement,
				" takes advantage of the defense's anticipation on the read,"+gainStatement,
				" surprises the defense on the read option and nimbly navigates the middle of the field for"+gainStatement,
				" surprises the defense on the read option and nimbly navigates the "+direction+" side of the field for"+gainStatement,
				" surprises the defense on the read option and elusively navigates the "+direction+" side of the field for"+gainStatement,
				" reads the defense on the option and elusively navigates the "+direction+" side of the field for"+gainStatement,
				" reads the defense on the option and hits the gas towards the "+direction+" side of the field for"+gainStatement,
				" reads the defense on the option and dashes towards the "+direction+" edge,"+gainStatement)
		case yards > 9:
			list = append(list, " opts to keep and finds a seam, powering ahead for ",
				" opts to keep and finds a gap on the "+direction+" side,"+gainStatement,
				" opts to keep and finds a seam on the "+direction+" side,"+gainStatement,
				" executes the option and finds a seam on the "+direction+" side,"+gainStatement,
				" executes the option and finds a gap on the "+direction+" side, powering through for ",
				" makes the right call on the read option, rushing for"+gainStatement)
		case yards > 4:
			list = append(list, " keeps on the read option and maneuvers,"+gainStatement,
				" reads the play well, keeping it,"+gainStatement,
				" reads the play well, executes the option,"+gainStatement,
				" surprises the defensive line on the option but is "+tackleVerb+" shortly on the "+direction+" side,"+gainStatement,
				" finds a gap on the defensive line on the option but is "+tackleVerb+" shortly on the "+direction+" side,"+gainStatement,
				" finds a seam on the defensive line on the option but is "+tackleVerb+" shortly on the "+direction+" side,"+gainStatement,
				" finds a seam through the read option but is stopped shortly on the "+direction+" side for"+gainStatement)
		case yards > 0:
			list = append(list, " decides to keep but is quickly "+tackleVerb+" after gaining ",
				" makes a quick read and pushes forward "+gainStatement,
				" makes a quick read and is "+tackleVerb+" quickly "+gainStatement,
				" makes a quick read and is "+tackleVerb+" shortly,"+gainStatement)
		case yards == 0:
			list = append(list, " is met at the line of scrimmage on the option, no gain on the play. ",
				" is immediately stopped at the line of scrimmage on the option, no gain on the play. ",
				" attempts to execute the read option but cannot trick the defense, no gain on the play. ",
				" attempts to execute the read option but is "+tackleVerb+" swiftly, no gain on the play. ",
				" attempts to execute the read option but is "+tackleVerb+" on the line of scrimmage, no gain on the play. ",
				" attempts to execute the read option but is shut down on the line of scrimmage, no gain on the play. ",
				" attempts the read option but is shut down on the line of scrimmage, no gain on the play. ",
				" attempts the read and is immediately "+tackleVerb+" on the line. No gain on the play. ",
				" is immediately "+tackleVerb+" by the defensive line. No gain on the play. ")
		default: // Negative yardage
			list = append(list,
				" makes a read but is "+tackleVerb+" behind the line,"+gainStatement,
				" makes a read but is met before the line of scrimmage,"+gainStatement,
				" makes a read but is "+tackleVerb+" swiftly,"+gainStatement,
				" attempts a read but is "+tackleVerb+" behind the line,"+gainStatement,
				" attempts a read but is met before the line of scrimmage,"+gainStatement,
				" opts to keep but is tackled for a loss of ",
			)
		}
	}
	return list
}

func getSpeedOptionList(yards int, isleft, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) []string {
	list := []string{}
	direction := "left"
	if isleft {
		direction = "right"
	}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	// Depending on the yards gained or lost, append appropriate descriptions
	if touchdown {
		list = append(list,
			" picks up the speed option flawlessly, dashing to the "+direction+" and into the end zone! TOUCHDOWN! ",
			" picks up the speed option flawlessly, dashing "+direction+" into the end zone! TOUCHDOWN! ",
			" takes the option and makes a mad dash to the "+direction+" into the end zone for a TOUCHDOWN! ",
			" takes the option and speeds off to the "+direction+" into the end zone for a TOUCHDOWN! ",
			" successfully picks up the option to the "+direction+", evades a tackle and makes it to the endzone! TOUCHDOWN! ",
			" sprints off to the "+direction+" on the option and breaks away -- he's heading for the endzone! It's a TOUCHDOWN! ",
			" breaks a tackle on the option play and breaks away to the "+direction+" side! He's going... GOING... we've got a TOUCHDOWN! ",
			" takes the option and speeds off down the "+direction+" sideline for a TOUCHDOWN! ",
			" takes the option and runs off down the "+direction+" sideline for a TOUCHDOWN! ",
			" takes the option and dashes down the "+direction+" sideline for a TOUCHDOWN! ",
			" takes the option and breaks away down the "+direction+" sideline... TOUCHDOWN! ")
	} else if outOfBounds {
		list = append(list,
			" takes the option and runner bolts to the "+direction+", stepping out after"+gainStatement,
			" takes the option and runner speeds to the "+direction+", stepping out after"+gainStatement,
			" takes the option and runner runs to the "+direction+", stepping out after"+gainStatement,
			" takes the option, heads towards the "+direction+" where he's pushed out after"+gainStatement,
			" takes the option, runs towards the "+direction+" where he's pushed out after"+gainStatement,
			" takes the option, speeds off towards the "+direction+" where he's pushed out after"+gainStatement,
		)
	} else if twoPtConversion && touchdown {
		list = append(list,
			" nails the speed option for a successful two-point conversion, with the runner breaking to the "+direction+" and crossing the plane! ",
			" capitalizes on the speed option, pitching it for a successful dash into the end zone for two points ! ")
	} else if twoPtConversion && !touchdown {
		list = append(list,
			" attempts the speed option for two points but the defense clamps down, thwarting the conversion to the "+direction+". ",
			" tries to convert with a speed option but the pitch to the "+direction+" is stopped, no extra points. ")
	} else if fumble {
		list = append(list,
			" makes a run for it on the speed option but loses control of the ball after gaining ",
			" makes a break for it on the speed option and he's lost the ball! It's a fumble! What a turn of events after ",
			" makes a break for it on the speed option and with the tackl- it's a fumble! What a turn of events after ",
			" makes a break for it on the speed option and with the tack- it's a fumble! What a turn of events after ",
			" makes a run for it on the option and loses the football on the tackle! ",
			" makes a run for it on the option and loses the ball after being brought down for ",
			" fumbles the pitch on the speed option, losing control of the ball after a minor gain of ",
			" fumbles the pitch on the speed option and loses the ball after a minor gain of ",
			" sees the speed option go awry as the ball hits the ground during the pitch, missing ",
			" mistakes the pitch and loses control of the ball, fumbling for ")
	} else if safety {
		list = append(list,
			" disastrously opts for a speed option that ends in a safety, with the runner tackled in the end zone! ",
			" tries a risky speed option in the end zone, which backfires for a safety against them! ",
			"attempts the speed option but cannot fool the defense, and that's a safety! ",
			"attempts the speed option but can't fool the defense, and that's a safety! ")
	} else {
		switch {
		case yards > 14:
			list = append(list,
				" tricks the defense on the speed option and tears down the "+direction+" for"+gainStatement,
				" tricks the defense on the speed option and runs to the "+direction+" for"+gainStatement,
				" surprises the defense on the speed option and runs to the "+direction+" for"+gainStatement,
				" surprises the defense on the speed option and tears down the "+direction+" for"+gainStatement,
				" masterfully executes the speed option and tears down the "+direction+" for"+gainStatement,
				" masterfully executes the speed option and runs to the "+direction+" for"+gainStatement,
				" breaks a tackle on the speed option and makes a run for the "+direction+" for"+gainStatement,
				" breaks a tackle on the speed option and makes a break for the "+direction+" side for"+gainStatement,
				" breaks a tackle on the speed option and dashes for the "+direction+" sideline for"+gainStatement,
				" breaks a tackle on the speed option and dashes downfield towards the "+direction+" sideline,"+gainStatement,
				" skillfully executes a speed option, leading to a sprint along the "+direction+" that racks up significant yardage of ")
		case yards > 9:
			list = append(list,
				" opts for the speed option, succeeding the pitch that results in a brisk "+direction+" run for"+gainStatement,
				" opts for the speed option, making a timely pitch that results in a swift "+direction+" run for"+gainStatement,
				" opts for the speed option, succeeding the pitch that results in a swift "+direction+" run for"+gainStatement,
				" eludes defenders on the speed option and makes some ground before being "+tackleVerb+","+gainStatement,
				" eludes defenders on the speed option and darts to the "+direction+" before being "+tackleVerb+","+gainStatement,
				" eludes defenders on the speed option and runs to the "+direction+" before being "+tackleVerb+","+gainStatement,
				" misses a tackle on the speed option and runs to the "+direction+" before being "+tackleVerb+","+gainStatement,
				" capitalizes on the speed option, breaking to the "+direction+" for"+gainStatement,
				" capitalizes on the option, breaking to the "+direction+" for"+gainStatement)
		case yards > 4:
			list = append(list,
				" carries out on the speed option to the "+direction+","+gainStatement,
				" goes with the speed option, going around the "+direction+" edge,"+gainStatement,
				" takes the speed option to the "+direction+" edge,"+gainStatement,
				" takes the speed option to the "+direction+" and makes some ground,"+gainStatement,
			)
		case yards > 0:
			list = append(list,
				" manages to eke out a few yards on a speed option to the "+direction+","+gainStatement,
				" manages to eke out a few yards on the option pitch to the "+direction+","+gainStatement,
				" edges out a few yards on the option pitch to the "+direction+","+gainStatement,
				" takes the option to the "+direction+" side and is "+tackleVerb+","+gainStatement,
				" sees a modest return on the speed option, with the pitch leading to"+gainStatement)
		case yards == 0:
			list = append(list, " is met at the line of scrimmage, no gain on the play. ",
				" is immediately stopped at the line of scrimmage, no gain on the play. ",
				" is immediately tackled by the defensive line. No gain on the play. ")
		default: // Negative yardage
			list = append(list,
				" is stopped on the speed option as two defenders spy both option players. Caught behind the line,"+gainStatement,
				" is stopped on the speed option as two defenders spy both option players,"+gainStatement,
				" has trouble with the option pitch, handles the ball and goes down behind the line,"+gainStatement,
				" has trouble with the option pitch and goes down behind the line,"+gainStatement,
				" hesitates on the speed option pitch and goes down behind the line,"+gainStatement,
				" hesitates on the speed option pitch and is "+tackleVerb+" behind the line,"+gainStatement,
				" can't escape the edge on the option and is "+tackleVerb+" behind the line,"+gainStatement,
				" misreads the speed option and is "+tackleVerb+" behind the line,"+gainStatement,
				" misreads the speed option and goes down behind the line,"+gainStatement,
				" sees the speed option falter, with the runner caught behind the line,"+gainStatement,
				" watches the speed option play crumble, resulting in negative yardage of ")
		}
	}
	return list
}

func getInvertedOptionList(yards int, isleft, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) []string {
	list := []string{}
	direction := "left"
	if isleft {
		direction = "right"
	}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	// Depending on the yards gained or lost, append appropriate descriptions
	if touchdown {
		list = append(list,
			" masterfully executes the inverted option, slicing through the defense for a TOUCHDOWN! ",
			" takes the inverted path, finds a seam, and bursts into the end zone for a TOUCHDOWN! ",
			" orchestrates the veer option beautifully, diving into the end zone for the score! TOUCHDOWN! ")
	} else if outOfBounds {
		list = append(list,
			" opts for the inverted option, veering to the "+direction+" and stepping out with a gain. ",
			" takes the outside path on the veer option, securing yards before going out of bounds. ")
	} else if twoPtConversion && touchdown {
		list = append(list,
			" converts the inverted option into two points, brilliantly finding the end zone! ",
			" capitalizes on the veer option, zipping into the end zone for the conversion! ")
	} else if twoPtConversion && !touchdown {
		list = append(list,
			" attempts the inverted option for the conversion but is halted just shy of success! ",
			" tries to punch in with the veer option but can't quite make it, no conversion. ")
	} else if fumble {
		list = append(list,
			" loses the handle while executing the inverted option, fumbling the ball away. ",
			" tries to maneuver on the veer option but coughs up the ball, turning it over. ")
	} else if safety {
		list = append(list,
			" is trapped in the end zone attempting the inverted option, leading to a safety! ",
			" can't escape the clutches on the veer option and is tackled for a safety. ")
	} else {
		switch {
		case yards > 14:
			list = append(list,
				" exploits the inverted option for"+gainStatement,
				" dazzles with a veer to the "+direction+", breaking free for"+gainStatement,
				" takes the veer to the "+direction+", running down the sideline before being"+tackleVerb+","+gainStatement,
			)
		case yards > 9:
			list = append(list,
				" drives the inverted option upfield, deftly navigating"+gainStatement,
				" veers to the "+direction+","+gainStatement,
				" maneuvers on the veer, skirting defenders and picking up ")
		case yards > 4:
			list = append(list,
				" grinds out yards with the inverted option,"+gainStatement,
				" grinds out yards with the inverted option before being "+tackleVerb+","+gainStatement,
				" finds a bit of daylight on the veer,"+gainStatement)
		case yards > 0:
			list = append(list,
				" squeezes a few yards out of the inverted option,"+gainStatement,
				" makes a marginal gain with the veer,"+gainStatement,
				" makes a marginal gain with the veer for"+gainStatement)
		case yards == 0:
			list = append(list, " is met at the line of scrimmage, no gain on the play. ",
				" is immediately stopped on the veer at the line of scrimmage, no gain on the play. ",
				" takes the inverted option and is immediately tackled by the defensive line. No gain on the play. ")
		default: // Negative yardage
			list = append(list, " is "+tackleVerb+" on the inverted option,"+gainStatement,
				" is "+tackleVerb+" on the inverted option behind the LOS,"+gainStatement,
				" is brought down on the inverted option behind the LOS for a loss of ",
				" attempts the inverted option but is "+tackleVerb+" behind the line,"+gainStatement,
				" attempts the inverted option but is stiffled behind the line,"+gainStatement,
				" encounters resistance on the veer, getting pushed back for ")
		}
	}
	return list
}

func getTripleOptionList(yards int, isleft, touchdown, outOfBounds, twoPtConversion, fumble, safety bool) []string {
	list := []string{}
	direction := "left"
	if isleft {
		direction = "right"
	}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	// Depending on the yards gained or lost, append appropriate descriptions
	if touchdown {
		list = append(list,
			" executes the triple option to perfection, resulting in a dynamic TOUCHDOWN run! ",
			" masterfully handles the triple option, finding the end zone with an impressive dash! TOUCHDOWN! ",
			" navigates the triple option with ease, leading to a sensational TOUCHDOWN! ")
	} else if outOfBounds {
		list = append(list,
			" opts to keep it on the triple option, darting to the "+direction+" and out of bounds for"+gainStatement,
			" pitches on the triple option and the runner scoots out of bounds after"+gainStatement)
	} else if twoPtConversion && touchdown {
		list = append(list,
			" brilliantly converts the triple option into a two-point success, crossing the goal line! ",
			" leverages the triple option for a crucial two-point conversion, punching it in! ")
	} else if twoPtConversion && !touchdown {
		list = append(list,
			" attempts a triple option for the conversion but the defense holds firm, denying the points. ",
			" runs the triple option on the two-point attempt but fails to find the end zone. ")
	} else if fumble {
		list = append(list,
			" loses control during the triple option play, resulting in a turnover. ",
			" fumbles the ball while executing the triple option, squandering the drive. ")
	} else if safety {
		list = append(list,
			" gets trapped for a safety while trying to execute the triple option in the end zone. ",
			" makes a risky decision on the triple option, leading to a safety. ")
	} else {
		switch {
		case yards > 14:
			list = append(list,
				" dismantles the defense with a well-executed triple option, tearing off"+gainStatement,
				" showcases exceptional decision-making on the triple option, sprinting for"+gainStatement,
				" takes off on the triple option, running down the "+direction+" sideline before being "+tackleVerb+"for"+gainStatement,
				" breaks free on the triple option and makes some distance towards the "+direction+" sideline before being "+tackleVerb+"for"+gainStatement,
				" executes the triple option takes off towards the "+direction+" sideline, evading defenders before being "+tackleVerb+"for"+gainStatement)
		case yards > 9:
			list = append(list,
				" skillfully maneuvers through the defense on the triple option,"+gainStatement,
				" breaks past the line on the triple option for"+gainStatement,
				" breaks past the line on the triple option before being "+tackleVerb+" for"+gainStatement,
				" finds the right lane on the triple option, makes some headway"+gainStatement,
				" surprises the defense on the triple option, marching down for"+gainStatement)
		case yards > 4:
			list = append(list,
				" gains positive yardage on the triple option, pushing forward to"+gainStatement,
				" takes off on the triple option before getting "+tackleVerb+","+gainStatement)
		case yards > 0:
			list = append(list,
				" ekes out some yardage on the triple option,"+gainStatement,
				" runs the triple option before being "+tackleVerb+","+gainStatement,
				" is "+tackleVerb+" after running the triple option,"+gainStatement,
				" advances the ball with a cautious triple option play,"+gainStatement)
		case yards == 0:
			list = append(list, " is met at the line of scrimmage, no gain on the play. ",
				" is immediately stopped at the line of scrimmage, no gain on the play. ",
				" is immediately tackled by the defensive line. No gain on the play. ")
		default: // Negative yardage
			list = append(list, " is quickly "+tackleVerb+" on the triple option, resulting in"+gainStatement,
				" has trouble with the triple option and is "+tackleVerb+" for"+gainStatement,
				" mishandles the triple option and is "+tackleVerb+" for"+gainStatement,
				" struggles to find space on the triple option, leading to"+gainStatement)
		}
	}
	return list
}

func getTackledVerb() string {
	list := []string{"brought down", "tackled", "taken down", "stopped", "manhandled", "swarmed"}
	return PickFromStringList(list)
}

func getGainSuffix(isGain bool, yards int) string {
	if isGain {
		list := []string{" a gain of ", " gaining ", " a moderate gain of ", " advancing for ", " a positive gain of ", " banking ", " racking up ", " picking up ", " adding "}
		if yards > 14 {
			list = append(list, " a massive gain of ", " a fantastic gain of ",
				" an incredible gain of ", " a massive gain of ", " massively gaining ", " a tremendous run of ", " a fantastic gain of ", " an incredible gain of ", " weaving through traffic for ", " a significant gain of ", " breaking free for ")
		}
		if yards > 9 {
			list = append(list, " a strong gain of ", " a strong gain of ", " strongly gaining ", " a pickup of ", " picking up ", " a solid gain of ", " a solid run of ", " amassing ", " a hard earned ", " powering ahead for ", " powering through for ")
		} else if yards > 4 {
			list = append(list, " a moderate gain of ", " moderately gaining ", " churning for ", " weaving through traffic for ",
				" trucking ahead for ", " a decent advance of ", " a respectable gain of ", " grinding out ", " a healthy gain of ", " bolting for ", " securing ", " notching up ", " bagging ", " collecting ")
		} else {
			list = append(list, " a short gain of ", " a short gain of ", " inching for ", " a small gain of ", " a minor gain of ",
				" a slight gain of ", " every inch fought for ", " adding ", " snagging ", " securing ", " only managing ", " a minimal gain of ", " bagging ", " eking out ", " scraping together ")
		}
		return PickFromStringList(list)
	} else {
		list := []string{" a loss of ", " getting pushed back for ", " losing ", " a loss of ", " a negative gain of "}
		return PickFromStringList(list)
	}
}
