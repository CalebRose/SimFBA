package util

import (
	"math"
	"strconv"
	"strings"
)

func GetPassStatement(yards int, offensiveFormation, playName, poa, recLabel string,
	touchdown, outOfBounds, twoPtConversion, fumble,
	safety, scramble, sack, complete, interception bool,
	turnOverLabel string) string {
	snapText := getSnapText(offensiveFormation)
	scrambleText := ""
	if scramble {
		scrambleText = getScrambleText(yards, touchdown)
		return snapText + scrambleText
	}
	if sack {
		sackText := getSackText(safety, fumble, touchdown)
		return snapText + sackText
	}
	throwStatement := getThrowStatement(yards, recLabel)
	if interception {
		intText := getInterceptText(yards, recLabel, turnOverLabel, fumble, touchdown)
		return snapText + throwStatement + intText
	}
	resultText := ""
	if !complete && len(recLabel) == 0 {
		resultText = getIncompleteThrowText(recLabel)
		return snapText + resultText
	}
	if complete {
		resultText = getCompleteThrowText(yards, recLabel, turnOverLabel, fumble, touchdown, twoPtConversion, outOfBounds, safety)
	} else {
		resultText = getIncompleteThrowText(recLabel)
	}

	finalString := snapText + scrambleText + throwStatement + resultText
	return finalString
}

func getIncompleteThrowText(recLabel string) string {
	var list []string
	if len(recLabel) > 0 {
		// If there is an intended receiver
		list = append(list, " and misses his man. Incomplete pass intended for "+recLabel,
			" and the ball is tipped off "+recLabel+"'s hands and onto the field. Incomplete pass intended for "+recLabel,
			" and the ball bounces off "+recLabel+"'s hands and it's incomplete. Pass intended for "+recLabel,
			" aims for "+recLabel+", but the pass falls incomplete under tight coverage.  Pass intended for "+recLabel+". ",
			" but the pass sails over his head. Incomplete. Pass intended for "+recLabel+". ",
			" but the throw is just out of reach. Incomplete. Pass intended for "+recLabel,
			" but the throw goes right through his hands and hits the turf. Incomplete. Pass intended for "+recLabel,
			" but the throw veers too far left. Incomplete. Pass intended for "+recLabel,
			" tries to connect with "+recLabel+", but the throw veers too far right. Incomplete. ",
			" but the pass is deflected! Incomplete. Pass intended for "+recLabel+". ",
			" but the pass is broken up! Incomplete.  Pass intended for "+recLabel+". ",
			" and it's overthrown, beyond the reach of "+recLabel+". Incomplete. ",
			" and the ball is too low and hits the turf. Incomplete pass intended for "+recLabel+". ",
		)
	} else {
		// If there is no open receiver
		list = append(list, " tries to find an open receiver but faces pressure. Throws it out of bounds. ",
			" He's flushed out of the pocket and has to throw it away. Incomplete. ",
			" He's forced out of the pocket and has to throw it away. Incomplete. ",
			" He's forced out of the pocket and has to throw it to the sideline. Incomplete. ",
			" He's forced out of the pocket and throws it out of bounds. Incomplete. ",
			" Cannot find an open man and chucks it out of bounds to avoid the sack. ",
			" He's in trouble and tosses the ball into the sidelines. Incomplete. ",
			" Feels the heat from the defense and throws the ball to the sideline. Incomplete. ",
		)
	}
	return PickFromStringList(list)
}

func getCompleteThrowText(yards int, recLabel, turnoverLabel string, fumble, touchdown, twoPtConversion, outofbounds, safety bool) string {
	baseList := []string{" and it's caught! ", " and he catches it! "}
	tackleVerb := getTackledVerb()
	gainStatement := getGainSuffix(yards > 0, yards)
	var list []string
	if fumble && !touchdown && !twoPtConversion {
		// Fumble & recovery
		fumex := getFumbleExpression()
		list = append(list, recLabel+" tries to get control of the ball and it fumbles loose! "+fumex+" ",
			recLabel+" can't seem to have control of the ball and loses it! "+fumex+" ",
			recLabel+" is brought down after the catch and the ball is loose! "+fumex+" ",
			recLabel+" makes the catch but gets hit hard and coughs up the football! "+fumex+" ",
			recLabel+" secures the pass but loses his grip when tackled! "+fumex+" ",
			recLabel+" hauls it in but the defender strips the ball away! "+fumex+" ",
			recLabel+" makes the reception and immediately gets leveled - the ball pops out! "+fumex+" ",
			recLabel+" catches it but fumbles on contact with the defender! "+fumex+" ",
			recLabel+" grabs the pass but can't hold onto it through the hit! "+fumex+" ",
			recLabel+" makes the catch but the ball is jarred loose by a big tackle! "+fumex+" ",
			recLabel+" secures the ball momentarily before it's knocked free! "+fumex+" ",
			recLabel+" gets his hands on it but loses control when he hits the ground! "+fumex+" ",
			recLabel+" makes the grab but the ball slips out during the tackle! "+fumex+" ",
			recLabel+" catches it cleanly but fumbles when multiple defenders converge! "+fumex+" ",
			recLabel+" pulls it in but the defender punches the ball out! "+fumex+" ",
			recLabel+" has it ripped away just as he secures the catch! "+fumex+" ",
		)
	} else if fumble && touchdown && !twoPtConversion {
		// Defensive return
		list = append(list, recLabel+" is quickly brought down - the ball is loose! The defense has recovered the ball and are taking it home! TOUCHDOWN! ",
			recLabel+" makes a run for it an- THE BALL IS STRIPPED FROM HIS HANDS! "+turnoverLabel+" makes a run for the other side! TOUCHDOWN! ",
			recLabel+" runs for it an- HE FUMBLES THE BALL! "+turnoverLabel+" scoops it up and is making a run for the endzone! No other man in sight! TOUCHDOWN! ",
			recLabel+" fights for control and it's stripped from his hands by "+turnoverLabel+"! "+turnoverLabel+" breaks away and makes marches down the field to the end zone! TOUCHDOWN! ",
			recLabel+" secures the catch but loses control! "+turnoverLabel+" pounces on the loose ball and races for the score! TOUCHDOWN! ",
			recLabel+" gets hit and the ball pops out! "+turnoverLabel+" scoops and scores in a flash! What a defensive touchdown! TOUCHDOWN! ",
			recLabel+" is gang tackled and the ball comes loose! "+turnoverLabel+" grabs it and weaves through traffic all the way to the house! TOUCHDOWN! ",
			recLabel+" struggles to maintain possession and it's jarred loose! "+turnoverLabel+" picks it up and nobody can catch him! TOUCHDOWN! ",
			recLabel+" takes a massive hit and fumbles! "+turnoverLabel+" recovers and turns on the jets for an incredible return! TOUCHDOWN! ",
			recLabel+" bobbles the catch and it falls to the turf! "+turnoverLabel+" alertly recovers and sprints to paydirt! TOUCHDOWN! ",
			recLabel+" has it ripped away by the defender! "+turnoverLabel+" immediately breaks free and outruns the entire offense! TOUCHDOWN! ",
			recLabel+" loses his grip after contact! "+turnoverLabel+" capitalizes with perfect field vision and takes it to the house! TOUCHDOWN! ",
			recLabel+" gets absolutely leveled and the ball flies out! "+turnoverLabel+" makes the recovery and shows incredible speed to the endzone! TOUCHDOWN! ",
			recLabel+" fights for extra yards but coughs it up! "+turnoverLabel+" makes a heads-up play and returns it for six! TOUCHDOWN! ",
			recLabel+" takes a crushing blow and loses the football! "+turnoverLabel+" snatches it up and burns everyone to the endzone! TOUCHDOWN! ",
		)
	} else if !fumble && touchdown && !twoPtConversion {
		// Deep throws for a touchdown
		if yards > 30 {
			list = append(list, recLabel+" runs it down the sideline, no DB in sight, and brings it the endzone! TOUCHDOWN! ",
				recLabel+" catches it out in the open and takes it straight home! TOUCHDOWN! ",
				recLabel+" catches it one-handed on a dive into the endzone! TOUCHDOWN! ",
				recLabel+" catches the deep ball and shakes off the coverage. He's going for the endzone! TOUCHDOWN! ",
				recLabel+" hauls in the deep pass and burns the secondary all the way to the house! TOUCHDOWN! ",
				recLabel+" makes an incredible catch over the defender and races to the endzone! TOUCHDOWN! ",
				recLabel+" breaks free on the deep route and makes the catch in stride for six! TOUCHDOWN! ",
				recLabel+" outmuscles the cornerback for the catch and sprints to paydirt! TOUCHDOWN! ",
				recLabel+" makes a spectacular diving catch and somehow keeps his feet for the score! TOUCHDOWN! ",
				recLabel+" catches the bomb and leaves the defense in the dust on his way to the endzone! TOUCHDOWN! ",
				recLabel+" makes an acrobatic catch at the goal line and tumbles into the endzone! TOUCHDOWN! ",
				recLabel+" beats double coverage and hauls in the deep ball for a stunning score! TOUCHDOWN! ",
				recLabel+" catches it over his shoulder and outpaces everyone to the endzone! TOUCHDOWN! ",
				recLabel+" makes the grab and jukes the safety for an incredible touchdown run! TOUCHDOWN! ",
				recLabel+" snags the pass in traffic and breaks multiple tackles on his way to glory! TOUCHDOWN! ",
			)
		} else if yards > 10 {
			// Medium throws for a touchdown
			list = append(list, recLabel+" evades a safety and runs right into the endzone! TOUCHDOWN! ",
				recLabel+" finds a seam in the defense, grabs the pass, and darts into the endzone! TOUCHDOWN! ",
				"With a quick cut, "+recLabel+" gets open, catches the pass, and zips into the endzone! TOUCHDOWN! ",
				recLabel+" latches onto the pass in stride and breaks the plane for a TOUCHDOWN! ",
				"In the red zone, "+recLabel+" secures the pass and sidesteps a defender to score! TOUCHDOWN! ",
				recLabel+" makes the catch and turns upfield, breaking a tackle on his way to the endzone! TOUCHDOWN! ",
				recLabel+" hauls in the pass at the 10 and races untouched to paydirt! TOUCHDOWN! ",
				recLabel+" hauls in the pass at the 9 and races untouched to paydirt! TOUCHDOWN! ",
				recLabel+" hauls in the pass at the 8 and races untouched to paydirt! TOUCHDOWN! ",
				recLabel+" hauls in the pass at the 7 and races untouched to paydirt! TOUCHDOWN! ",
				recLabel+" hauls in the pass at the 6 and races untouched to paydirt! TOUCHDOWN! ",
				recLabel+" hauls in the pass at the 5 and races untouched to paydirt! TOUCHDOWN! ",
				recLabel+" catches it on the comeback route and spins away from defenders into the endzone! TOUCHDOWN! ",
				recLabel+" snags the pass over the middle and fights through contact for the score! TOUCHDOWN! ",
				recLabel+" makes a fantastic catch and bulldozes his way into the endzone! TOUCHDOWN! ",
				recLabel+" secures the pass and makes a brilliant cut to the corner for six! TOUCHDOWN! ",
				recLabel+" grabs the throw in traffic and powers through multiple defenders to score! TOUCHDOWN! ",
				recLabel+" catches the pass at the goal line and dives across for the touchdown! TOUCHDOWN! ",
				recLabel+" makes the reception and uses a stiff arm to reach the endzone! TOUCHDOWN! ",
				recLabel+" pulls in the pass and makes a spectacular hurdle into the endzone! TOUCHDOWN! ",
				recLabel+" catches it cleanly and accelerates past the linebacker for the score! TOUCHDOWN! ",
				recLabel+" secures the ball and makes a defender miss with a perfect juke to score! TOUCHDOWN! ",
				recLabel+" makes the grab and shows great vision to find the endzone! TOUCHDOWN! ",
			)
		} else {
			// Short throws within the red zone/endzone
			list = append(list, recLabel+" makes the catch right into the endzone, TOUCHDOWN! ",
				recLabel+" catches it in tight coverage and steps in the endzone before being pushed out of bounds! TOUCHDOWN! ",
				recLabel+" runs it down and dives for the endzone! TOUCHDOWN! ",
				recLabel+" makes the catch on a dive into the endzone! TOUCHDOWN! ",
				recLabel+" makes the catch in traffic and is pushed out of bounds in the endzone. TOUCHDOWN! ",
				recLabel+" makes the catch out in the open and steps right into the endzone! TOUCHDOWN! ",
				recLabel+" jukes a defender and makes a mad dash right into the endzone! TOUCHDOWN! ",
				recLabel+" grabs the quick slant and barrels over the line for a TOUCHDOWN! ",
				"In a crowded endzone, "+recLabel+" snatches the pass for a decisive TOUCHDOWN! ",
				recLabel+" makes a pivotal catch in the corner of the endzone! TOUCHDOWN! ",
				"Under pressure, "+recLabel+" secures the pass and tumbles into the endzone! TOUCHDOWN! ",
				recLabel+" scoops up the low throw and turns it into six points! TOUCHDOWN! ",
				"Amidst tight coverage, "+recLabel+" reels in the pass and plants his feet in the endzone! TOUCHDOWN! ",
			)
		}
	} else if !fumble && touchdown && twoPtConversion {
		list = append(list, recLabel+" has caught it in the endzone and succeeds on the two point conversion! ",
			recLabel+" catches it in tight coverage and banks on the two point conversion! ",
			recLabel+" makes the catch in traffic and succeeds on the two point conversion! ",
			recLabel+" on the two point conversion in the endzone! ",
			recLabel+" secures the pass in the corner of the endzone for the two-point conversion! ",
			recLabel+" makes a diving catch to complete the two-point conversion! ",
			recLabel+" fights through contact to make the catch for two points! ",
			recLabel+" pulls in the fade route perfectly for the successful two-point try! ",
			recLabel+" makes a spectacular grab over the defender for two points! ",
			recLabel+" catches the low throw and gets both feet down for the conversion! ",
			recLabel+" beats his man on the slant and secures the two-point conversion! ",
			recLabel+" makes an acrobatic catch at the goal line for two points! ",
			recLabel+" hauls in the pass despite heavy pressure for the conversion! ",
			recLabel+" makes the tough catch in double coverage for two points! ",
			recLabel+" secures the ball with perfect timing for the two-point conversion! ",
			recLabel+" makes a leaping catch to convert the two-point attempt! ",
			recLabel+" catches it at the back of the endzone for the successful conversion! ",
			recLabel+" powers through the defender's contact to complete the two-point try! ",
		)
	} else if safety {
		// Safety
		list = append(list, recLabel+" struggles to make it out of the endzone and is "+tackleVerb+"! Safety! ",
			recLabel+" tries to find an open lane with the catch, but is "+tackleVerb+" in the endzone! Safety! ",
			recLabel+" is swarmed in the endzone and is brought down. Safety! ",
			recLabel+" makes the catch but has nowhere to go and is "+tackleVerb+" in the endzone! Safety! ",
			recLabel+" secures the pass but gets immediately surrounded and "+tackleVerb+"! Safety! ",
			recLabel+" catches it at the goal line but is pushed back into the endzone and "+tackleVerb+"! Safety! ",
			recLabel+" hauls in the pass but the defense collapses on him in the endzone! Safety! ",
			recLabel+" makes the reception but is overwhelmed by defenders in the endzone! Safety! ",
			recLabel+" grabs the ball but can't escape the endzone before being "+tackleVerb+"! Safety! ",
			recLabel+" catches it cleanly but gets gang tackled in the endzone! Safety! ",
			recLabel+" secures the pass but is immediately hit and "+tackleVerb+" for a safety! ",
			recLabel+" makes the catch but loses his footing and goes down in the endzone! Safety! ",
			recLabel+" pulls in the pass but the defense swarms him instantly! Safety! ",
			recLabel+" catches it but gets crushed by multiple defenders in the endzone! Safety! ",
			recLabel+" makes the grab but has no room to maneuver and is "+tackleVerb+"! Safety! ",
		)
	} else if outofbounds {
		list = append(list,
			recLabel+" makes the catch and steps out of bounds after "+gainStatement,
			recLabel+" hauls in the pass and quickly goes out of bounds, stopping the clock after"+gainStatement,
			recLabel+" grabs the throw and is immediately out of bounds,"+gainStatement,
			recLabel+" grabs the ball and is pushed out of bounds for"+gainStatement,
			recLabel+" with the catch and gets out of bounds for"+gainStatement,
			recLabel+" with the catch in traffic and steps out of bounds for"+gainStatement,
			recLabel+" catches the ball and tiptoes the sideline before stepping out,"+gainStatement,
			recLabel+" secures the pass near the boundary and wisely steps out to stop the clock,"+gainStatement,
			recLabel+" makes a spectacular sideline catch and is forced out by the defender,"+gainStatement,
			recLabel+" hauls in the pass and dives toward the sideline, landing out of bounds,"+gainStatement,
			recLabel+" catches it in bounds but his momentum carries him out,"+gainStatement,
			recLabel+" makes the reception and strategically heads for the sideline,"+gainStatement,
			recLabel+" grabs the pass and gets shoved out of bounds by multiple defenders,"+gainStatement,
			recLabel+" secures the ball and makes a smart decision to step out and preserve the clock,"+gainStatement,
			recLabel+" makes the catch under pressure and is driven out of bounds,"+gainStatement,
			recLabel+" pulls in the pass and gets knocked out of bounds by the defense,"+gainStatement,
			recLabel+" catches it cleanly and heads straight for the sideline,"+gainStatement,
			recLabel+" makes the grab and is immediately pushed toward the boundary,"+gainStatement,
			recLabel+" secures the pass and uses his awareness to get out of bounds,"+gainStatement,
			recLabel+" makes a toe-tapping catch right at the sideline,"+gainStatement,
			recLabel+" hauls it in and fights his way toward the boundary before stepping out,"+gainStatement,
			recLabel+" catches the pass and smartly protects the ball while going out,"+gainStatement,
			recLabel+" makes the reception and gets bumped out of bounds by contact,"+gainStatement,
			recLabel+" grabs the ball and is wrestled out of bounds by the defender,"+gainStatement,
			recLabel+" secures the catch and is tackled toward the sideline, stepping out,"+gainStatement)
	} else {
		tackleVerb := getTackledVerb()
		switch {
		case yards > 25:
			list = append(list,
				recLabel+" catches the deep ball and is "+tackleVerb+" after"+gainStatement,
				recLabel+" catches the ball and makes some significant headway before being "+tackleVerb+" after"+gainStatement,
				recLabel+" catches the ball and makes a run for it before being "+tackleVerb+" after"+gainStatement,
				recLabel+" catches the deep ball in traffic and is immediately "+tackleVerb+" after"+gainStatement,
				recLabel+" with the catch and makes a run through the secondary before being tackled "+tackleVerb+" after"+gainStatement,
				recLabel+" makes the catch out in the open and is chased down. Makes significant yardage with"+gainStatement,
				recLabel+" wows the stadium with an incredible catch and is "+tackleVerb+" by the defense. A great impact on the drive for"+gainStatement,
				recLabel+" hauls in the long pass and is "+tackleVerb+" after"+gainStatement,
				recLabel+" on the streak with the catch and is "+tackleVerb+" after"+gainStatement,
			)
		case yards > 14:
			list = append(list,
				recLabel+" snatches the pass out of the air and is "+tackleVerb+" after"+gainStatement,
				recLabel+" secures the catch and is "+tackleVerb+", but not before"+gainStatement,
				recLabel+" makes the catch in traffic and is "+tackleVerb+","+gainStatement,
				recLabel+" makes the catch in coverage and is "+tackleVerb+","+gainStatement,
				recLabel+" makes the catch nearly escapes but is "+tackleVerb+","+gainStatement,
				recLabel+" makes the catch nearly escapes but is "+tackleVerb+" for"+gainStatement,
				recLabel+" with the catch and is "+tackleVerb+","+gainStatement,
				recLabel+" with the catch and is "+tackleVerb+" for"+gainStatement,
				recLabel+" catches the ball and is "+tackleVerb+" for"+gainStatement,
				recLabel+" secures the pass and is "+tackleVerb+" for"+gainStatement,
				recLabel+" on the post with the catch and is "+tackleVerb+" after"+gainStatement,
			)
		case yards > 9:
			list = append(list,
				recLabel+" grabs the throw and after a quick move is "+tackleVerb+", marking "+gainStatement,
				recLabel+" makes the catch in coverage and is "+tackleVerb+","+gainStatement,
				recLabel+" with the catch and is "+tackleVerb+","+gainStatement,
				recLabel+" with the catch and is "+tackleVerb+" for"+gainStatement,
				recLabel+" catches the ball and is "+tackleVerb+" for"+gainStatement,
				recLabel+" secures the pass and is "+tackleVerb+" for"+gainStatement,
				recLabel+" on the flag route with the catch and is "+tackleVerb+" after"+gainStatement,
				recLabel+" catches the pass and fights his way through to add "+gainStatement+" before being "+tackleVerb)
		case yards > 4:
			list = append(list,
				recLabel+" makes a short catch and is "+tackleVerb+" after"+gainStatement,
				recLabel+" pulls in the pass and is "+tackleVerb+" after"+gainStatement,
				recLabel+" on the short route with the catch and is "+tackleVerb+" after"+gainStatement,
				recLabel+" on the slant with the catch and is "+tackleVerb+" after"+gainStatement,
				recLabel+" on the hook with the catch and is "+tackleVerb+" after"+gainStatement,
				recLabel+" secures the pass over the middle and is "+tackleVerb+" after"+gainStatement,
				recLabel+" makes the catch on the comeback route and is "+tackleVerb+" after"+gainStatement,
				recLabel+" hauls in the short pass and is immediately "+tackleVerb+" after"+gainStatement,
				recLabel+" catches the ball on the curl route and is "+tackleVerb+" after"+gainStatement,
				recLabel+" grabs the pass on the quick out and is "+tackleVerb+" after"+gainStatement,
				recLabel+" makes the reception on the drag route and is "+tackleVerb+" after"+gainStatement,
				recLabel+" secures the catch on the hitch and is "+tackleVerb+" after"+gainStatement,
				recLabel+" pulls down the pass on the crossing route and is "+tackleVerb+" after"+gainStatement,
				recLabel+" makes the grab on the quick slant and is "+tackleVerb+" after"+gainStatement,
				recLabel+" catches it cleanly on the short pattern and is "+tackleVerb+" after"+gainStatement,
				recLabel+" secures the ball on the underneath route and is "+tackleVerb+" after"+gainStatement,
				recLabel+" makes the catch in tight coverage and is "+tackleVerb+" after"+gainStatement,
				recLabel+" hauls in the pass on the timing route and is "+tackleVerb+" after"+gainStatement,
				recLabel+" grabs the throw on the quick game and is "+tackleVerb+" after"+gainStatement,
				recLabel+" catches the ball on the checkdown and is "+tackleVerb+" after"+gainStatement,
				recLabel+" secures the pass on the dig and is "+tackleVerb+" after"+gainStatement,
				recLabel+" makes the reception on the bubble screen and is "+tackleVerb+" after"+gainStatement,
				recLabel+" pulls in the pass on the shallow cross and is "+tackleVerb+" after"+gainStatement,
				recLabel+" catches it on the quick game route and is "+tackleVerb+" after"+gainStatement,
				recLabel+" makes the grab on the intermediate route and is "+tackleVerb+" after"+gainStatement,
			)
		case yards > 0:
			list = append(list,
				recLabel+" makes the catch and is quickly "+tackleVerb+" after"+gainStatement,
				recLabel+" makes a quick catch and is "+tackleVerb+" after"+gainStatement,
				recLabel+" on the quick route and is "+tackleVerb+" after"+gainStatement,
				recLabel+" on the slants, faces coverage and is "+tackleVerb+" after"+gainStatement,
				recLabel+" on the slants, faces coverage and is "+tackleVerb+" for"+gainStatement,
				recLabel+" barely gets the catch and is immediately "+tackleVerb+","+gainStatement,
				recLabel+" barely makes the catch and is quickly "+tackleVerb+","+gainStatement,
				recLabel+" barely makes the catch and is quickly "+tackleVerb+" after"+gainStatement,
				recLabel+" barely makes the catch and is quickly "+tackleVerb+" for"+gainStatement,
				recLabel+" makes the reception and is "+tackleVerb+" with"+gainStatement)
		case yards == 0:
			list = append(list,
				recLabel+" catches the ball but is "+tackleVerb+" at the spot, no gain on the play. ",
				recLabel+" catches the ball on the line but is "+tackleVerb+", no gain on the play. ",
				recLabel+" makes the catch but is "+tackleVerb+" on the line, no gain on the play. ",
				recLabel+" makes the catch but is "+tackleVerb+" on the line, no progress made. ",
				recLabel+" secures the pass but is instantly "+tackleVerb+", no progress made. ",
				recLabel+" hauls in the pass but is immediately "+tackleVerb+" at the line of scrimmage. ",
				recLabel+" makes the reception but can't advance, "+tackleVerb+" right where he caught it. ",
				recLabel+" catches it cleanly but is swarmed by defenders at the spot, no gain. ",
				recLabel+" secures the ball but faces a wall of defenders, "+tackleVerb+" for no gain. ",
				recLabel+" pulls in the pass but is met immediately by the defense, no advancement. ",
				recLabel+" makes the catch but is wrapped up instantly, "+tackleVerb+" at the original spot. ",
				recLabel+" grabs the pass but multiple defenders converge, "+tackleVerb+" for no gain. ",
				recLabel+" catches the ball in traffic and is "+tackleVerb+" right at the catch point. ",
				recLabel+" secures the reception but the defense closes in fast, no forward progress. ",
				recLabel+" makes the grab but is hit immediately upon catching it, "+tackleVerb+" for no gain. ",
				recLabel+" hauls it in but faces immediate contact, "+tackleVerb+" at the line. ",
				recLabel+" catches the pass but is engulfed by defenders, no room to move forward. ",
				recLabel+" makes the reception but is gang tackled at the spot, no gain on the play. ",
				recLabel+" secures the ball but the defense reacts quickly, "+tackleVerb+" for no advancement. ",
				recLabel+" pulls down the pass but faces immediate pressure, "+tackleVerb+" where he stands. ",
				recLabel+" catches it but is met by multiple defenders, "+tackleVerb+" at the original line. ",
				recLabel+" makes the catch but the defense swarms him instantly, no forward progress made. ")

		default:
			list = append(list,
				recLabel+" is "+tackleVerb+" for a loss after catching the ball, a tough break,"+gainStatement,
				recLabel+" catches it but is "+tackleVerb+" behind the line,"+gainStatement,
				recLabel+" tries to make some headway but is pushed back, a tough break for"+gainStatement,
				recLabel+" makes the catch but is quickly swarmed for a loss,"+gainStatement,
				recLabel+" secures the pass but is driven backward by the defense,"+gainStatement,
				recLabel+" hauls in the reception but gets hit hard and pushed back,"+gainStatement,
				recLabel+" makes the catch but faces immediate pressure and loses ground,"+gainStatement,
				recLabel+" catches the ball but is overwhelmed by defenders and pushed back,"+gainStatement,
				recLabel+" grabs the pass but is gang tackled and driven backward,"+gainStatement,
				recLabel+" secures the reception but the defense collapses on him for a loss,"+gainStatement,
				recLabel+" makes the catch but is met by multiple defenders and loses yardage,"+gainStatement,
				recLabel+" hauls it in but gets crushed by the defense and pushed behind the line,"+gainStatement,
				recLabel+" catches the pass but is immediately surrounded and driven back,"+gainStatement,
				recLabel+" makes the reception but faces a swarm of defenders and loses ground,"+gainStatement,
				recLabel+" secures the ball but is overpowered by the defense and pushed back,"+gainStatement,
				recLabel+" pulls in the pass but gets hit by multiple tacklers and loses yardage,"+gainStatement,
				recLabel+" catches it cleanly but is wrapped up and driven backward,"+gainStatement,
				recLabel+" makes the grab but the defensive pursuit catches up and pushes him back,"+gainStatement,
				recLabel+" hauls in the pass but is leveled by defenders and loses ground,"+gainStatement,
				recLabel+" secures the reception but gets blindsided and knocked backward,"+gainStatement,
				recLabel+" catches the ball but is met with fierce contact and driven back,"+gainStatement,
				recLabel+" makes the catch but is immediately hit and loses his footing, falling backward,"+gainStatement,
				recLabel+" grabs the pass but is tackled aggressively and pushed well behind the line,"+gainStatement,
				recLabel+" secures the ball but the defense swarms him and drives him backward,"+gainStatement,
				recLabel+" makes the reception but loses his balance and stumbles backward,"+gainStatement,
			)
		}
	}
	return PickFromStringList(baseList) + PickFromStringList(list)
}

func getInterceptText(yards int, recLabel, turnOverLabel string, fumble, touchdown bool) string {
	absYards := math.Abs(float64(yards))
	yardsInt := int(absYards)
	ydStr := strconv.Itoa(yardsInt)
	yardsStr := GetYardsString(int8(yardsInt))
	intVerb := getInterceptVerb()
	var list []string
	// Very rare case -- the problem is that I don't think we have the capacity to tell who scored based on the play data tangible
	if fumble && touchdown {
		list = append(list, " and he's "+intVerb+"! Caught by "+turnOverLabel+" with the catch and he makes a run for it! Brought down, an- the ball is lose! It looks like it's a fight for it, and it's picked up! The player's making it to the endzone! TOUCHDOWN! ",
			" and he's "+intVerb+"! "+turnOverLabel+" who dashes forward, and loses the ball! And it's scooped up! A crazy turn of events, as the ball is NOW being brought to the endzone! TOUCHDOWN! ")
	} else if fumble && !touchdown {
		fumex := getFumbleExpression()
		list = append(list, " and he's "+intVerb+"! Caught by "+turnOverLabel+" with the catch and he makes a run for it! Brought down, an- the ball is lose! "+fumex+" ",
			" and he's "+intVerb+"! "+turnOverLabel+" tries to make a break for it, but fumbles the football! "+fumex+" ",
			" and he's "+intVerb+"! "+turnOverLabel+" tries to make a break for it, but fumbles the football! "+fumex+" ")
	} else if !fumble && touchdown {
		list = append(list, " and he's "+intVerb+"! Caught by "+turnOverLabel+" with the catch and he makes a return all the way into the endzone! TOUCHDOWN! ",
			" and he's "+intVerb+"! Caught by "+turnOverLabel+" who takes it all the way back! TOUCHDOWN! What an incredible return. ",
			" and he's "+intVerb+"! Caught by "+turnOverLabel+", now making a fantastic return to all the way back! TOUCHDOWN! ",
			" and he's "+intVerb+"! "+turnOverLabel+" with the catch and he's breaking away from the pack, this is going to be a return all the way for a TOUCHDOWN! ",
			" and he's "+intVerb+"! "+turnOverLabel+" turns it over and is taking it all the way back! No one can stop him! TOUCHDOWN! ",
			" and he's "+intVerb+"! "+turnOverLabel+" turns it over and is taking it all the way back! The fan's are in shock and awe! TOUCHDOWN! ",
			" and he's "+intVerb+"! "+turnOverLabel+" turns it over and there is no one to stop him! He's going all the way to the endzone! What a turn of events, TOUCHDOWN! ",
			" and he's "+intVerb+"! "+turnOverLabel+" turns it over and there is no one to stop him! He's going all the way to the endzone! What a turn of events, TOUCHDOWN! ",
			" and he's "+intVerb+"! There it is! "+turnOverLabel+" is going towards the endzone with no one to stop him. A pick six! TOUCHDOWN! ",
			" and he's "+intVerb+"! "+turnOverLabel+" is going towards the endzone with no one to stop him. A pick six! TOUCHDOWN! ",
			" and it flies into the hands of "+turnOverLabel+"! "+turnOverLabel+" marches down the field with no one else in sight! TOUCHDOWN! ")
	} else {
		baseOptions := []string{" and he's " + intVerb + "! Caught by " + turnOverLabel + " with the catch and he makes a return for " + ydStr + yardsStr,
			" and it flies into the hands of " + turnOverLabel + "! He manages to get " + ydStr + " yards before being taken down. ",
			" and he's intercepted! A one-handed catch by " + turnOverLabel + "! And just like that he's brought down after " + ydStr + yardsStr,
			" and the ball is tipped off - it's " + intVerb + "! " + turnOverLabel + " with the catch and a return for " + ydStr + yardsStr,
			" and the ball is tipped off - it's " + intVerb + "! " + turnOverLabel + " catches it and manages a return for " + ydStr + yardsStr,
			" and he's " + intVerb + "! " + turnOverLabel + " with the catch and runs it down for " + ydStr + yardsStr,
			" and he's " + intVerb + "! " + turnOverLabel + " with the catch and runs it down for " + ydStr + yardsStr,
			" - OH and he's " + intVerb + "! " + turnOverLabel + " with the catch! Returns it downfield for " + ydStr + yardsStr,
			" and he's " + intVerb + "! " + turnOverLabel + " with the catch and runs it down for " + ydStr + yardsStr}
		base := PickFromStringList(baseOptions)
		if yards > 19 {
			list = append(list, base+"An impressive return following the pick! ",
				base+"Great coverage by "+turnOverLabel+" on "+recLabel+", and an amazing return too. ",
				base+"Great coverage by "+turnOverLabel+" on "+recLabel+", seizing the opportunity for value field position too. ")
		} else if yards > 9 {
			list = append(list, base+"A solid effort on the interception! ",
				base+turnOverLabel+" capitalized on that play and made a fantastic play for the defense. ")
		} else if yards > 4 {
			list = append(list, base+"Every yard counts as they fight to gain ground post-interception. ",
				base+"Fantastic coverage by "+turnOverLabel+" and a crucial turnover as well. ",
				base+turnOverLabel+" didn't lose sight on "+recLabel+" and seized this fantastic opportunity. ")
		} else {
			list = append(list, base+"Not much room to run, but the turnover is crucial. ",
				base+"Not much room to run, but what a crucial turnover! ",
				base+"Great coverage on "+recLabel+" and a turnover that could shake things up! ",
				base+turnOverLabel+" didn't lose sight on "+recLabel+", kept him covered, and seized this fantastic opportunity. ")
		}
	}
	return PickFromStringList(list)
}

func getSackText(safety, fumble, touchdown bool) string {
	var list []string
	if safety {
		list = append(list,
			"Is sacked in the end zone, resulting in a safety! ",
			"Can't escape the grasp of the defenders, leading to a safety with that sack! ",
			"Is taken down in the end zone for a safety, a crucial play! ",
			"The pressure gets to him, resulting in a sack for a safety! ",
			"Is overwhelmed by the defense and sacked for a safety! ")
	} else if fumble && !touchdown {
		fumex := getFumbleExpression()
		list = append(list,
			"Tries to evade the rush but is sacked! What's this? The ball is loose! "+fumex+" ",
			"Takes too long to find a man and the defense has broken through! A sack on the play -- and the ball has fumbled! "+fumex+" ",
			"Tries to throw it away but is sacked on the play! And wi- the ball is loose onto the field! "+fumex+" ",
			"Can't find a man and is sacked AND fumbles the ball! "+fumex+" ",
			"Can't evade the rush and is sack. What's this? He's lost his grip while being sacked, and the ball is loose on the field! ",
			"And he's brought down by the pass rush. The hit caused a fumble! The quarterback loses the ball as he's sacked! ",
			"The pocket collapses, and he's sacked! What's this? The ball is knocked loose! "+fumex+" ",
			"The defense has broken through and he's taken down hard - and the ball pops out! A fumble during the sack! "+fumex+" ",
		)
	} else if fumble && touchdown {
		list = append(list,
			"Can't find a man and is sacked on the play! What's this? There's a fumble on the field! The defense has scooped up the ball and is going... all the way! TOUCHDOWN! ",
			"The defense has broken through has sacked the quarterbac- and the ball is fumbled! It's been scooped up by the defense and is being returned. Not a man in sight -- that's a TOUCHDOWN! ",
			"Hesitates on a throw and is sacked. An- what's this? The ball is loose! It's a disaster for the offense as the defense has scooped up the ball and it's been taken back for a TOUCHDOWN! ",
			"The pocket collapses, and he's sacked! What's this? The ball has fumbled on the field! The defense capitalizes on the fumble with the recovery and a return to the endzone - TOUCHDOWN! ",
			"Takes too long on the throw and is sacked - the ball comes loose! The defense has scooped it up and is making a return for the endzone! TOUCHDOWN! ",
		)
	} else {
		list = append(list, "Hesitates on throwing the ball and is sacked on the play! ",
			"Tries to look for an open man, but is brought down by the defense. A huge sack on the play! ",
			"Can't find an open receiver and is sacked behind the line! ",
			"The pocket collapses, and he's sacked! ",
			"The pocket collapses, and he can't scramble out. A sack on the play! ",
			"What's this? The rush has overwhelmed the offensive line and the QB's taken down for a sack! ",
			"Tries to evade the rush but is sacked! ",
			"Takes too long to find a man and the defense has broken through! A sack on the play! ",
			"Tries to throw it away but is sacked on the play! ",
			"Is wrapped up and sacked, a significant loss on the play! ",
			"Faces a fierce pass rush and is sacked, thwarting the drive! ",
			"The defense breaks through and he's sacked, a big play! ",
			"Attempts to scramble but is caught and sacked! ",
			"Holds onto the ball too long and is sacked by the oncoming defenders! ",
			"Is hit and sacked, the defensive line breaking through! ",
		)
	}
	return PickFromStringList(list)
}

func getScrambleText(yards int, touchdown bool) string {
	gainStatement := getGainSuffix(yards > 0, yards)
	yardsStr := GetYardsString(int8(yards))
	if !touchdown {
		list := []string{"Leaves the pocket on and scrambles for " + gainStatement,
			"Sees pressure and scrambles out of the pocket for" + gainStatement + yardsStr,
			"Scrambles out of the pocket for" + gainStatement + yardsStr,
			"Leaves the pocket on a scramble, evading defenders for" + gainStatement + yardsStr,
			"Tucks the ball and scrambles, looking for open space. Runs for" + gainStatement + yardsStr,
			"Dodges the rush and scrambles out to the side for" + gainStatement + yardsStr,
			"Takes evasive action and scrambles to avoid a sack," + gainStatement + yardsStr,
			"Finds no one open and decides to scramble for yardage," + gainStatement + yardsStr,
			"Breaks away from pressure, scrambles for" + gainStatement + yardsStr,
			"Under duress, elects to scramble out of the pocket for" + gainStatement + yardsStr,
			"Sees the pocket collapsing and takes off on a scramble. Running for" + gainStatement + yardsStr,
			"Sees an opening and quickly scrambles out of the pocket for" + gainStatement + yardsStr,
			"Avoids the sack with a quick scramble to the side," + gainStatement + yardsStr,
			"Uses his legs to escape the pocket, runs for" + gainStatement + yardsStr,
			"Finds a lane and scrambles to exploit the gap, runs for" + gainStatement + yardsStr,
		}
		return PickFromStringList(list)
	}
	list := []string{"Leaves the pocket on and scrambles into the endzone for the TOUCHDOWN! ",
		"Sees pressure and scrambles out of the pocket gets right into the endzone! TOUCHDOWN! ",
		"Scrambles out of the pocket and makes a dive right into the endzone! TOUCHDOWN! ",
		"Leaves the pocket on a scramble, evading defenders on the goalline and makes it into the endzone - TOUCHDOWN! ",
		"Tucks the ball and scrambles for an open gap on the goalline. Dives into the endzone for the TOUCHDOWN! ",
		"Takes evasive action and scrambles to avoid a sack, runs to the edge of the goalline and makes it in - TOUCHDOWN! ",
		"Finds no one open and decides to scramble for the endzone. Succeeds, it's a TOUCHDOWN! ",
		"Sees the pocket collapsing and takes off on a scramble. Runs for the goalline and makes it into the endzone - TOUCHDOWN! ",
		"Avoids the sack with a quick scramble right into the endzone - TOUCHDOWN! ",
	}
	return PickFromStringList(list)
}

func getSnapText(form string) string {
	list := []string{" takes the snap. "}
	isShotgunPlay := CheckSubstring(form, "Gun")
	if isShotgunPlay {
		list = append(list, " drops back to pass. ",
			" takes the snap in the shotgun formation. ",
			" gets the ball in the shotgun. ",
			" fields the snap in the shotgun. ",
			" catches the snap while in the shotgun, scanning for options. ",
			" from the shotgun, secures the snap and setups the throw. ",
		)
	} else {
		list = append(list, " takes the snap from under center. ",
			" takes the ball from under center. ",
			" gets the snap and drops back. ",
			" receives the snap under center. Looks to pass. ",
			" under center, takes the snap and looks to pass. ",
			" collects the snap under center and prepares to throw. ",
			" snaps up the ball from under center and eyes his targets. ")
	}

	return PickFromStringList(list)
}

func getFumbleExpression() string {
	list := []string{
		"It's a fight for the ball!",
		"Both sides are scrambling to recover the ball.",
		"It's a scramble to recover the football!",
		"Everyone's scrambling for the ball in an attempt to recover it.",
		"The defense is scrambling to recover it!",
		"A potential turnover here!",
		"Both sides are trying to recover!",
		"A scramble for the ball has ensued!",
		"It's a fight for the pigskin!",
	}

	return PickFromStringList(list)
}

func getThrowStatement(yards int, recLabel string) string {
	throwVerb := getThrowingVerb(yards)
	distance := getDistance(yards)
	list := []string{throwVerb + " it ",
		throwVerb + " it to " + recLabel,
		throwVerb + " it " + distance + " to " + recLabel,
		throwVerb + " it " + distance + ", targeting " + recLabel,
		throwVerb + " the ball to " + recLabel,
		throwVerb + " the pass " + distance + " to " + recLabel,
		throwVerb + " a pass to " + recLabel,
		throwVerb + " the football " + distance,
		throwVerb + " it in the direction of " + recLabel,
		throwVerb + " the ball " + distance + ", looking for " + recLabel,
		throwVerb + " it toward " + recLabel,
		throwVerb + " the pass in the direction of " + recLabel,
		throwVerb + " a strike to " + recLabel,
		throwVerb + " it " + distance + " in search of " + recLabel,
		throwVerb + " the football to " + recLabel + " " + distance,
		throwVerb + " it " + distance + " where " + recLabel + " is running",
		throwVerb + " the ball in " + recLabel + "'s direction",
		throwVerb + " it to his target, " + recLabel,
		throwVerb + " a pass " + distance + " for " + recLabel}
	return PickFromStringList(list)
}

func getThrowingVerb(yards int) string {
	list := []string{"Throws", "Slings", "Passes", "Fires", "Lobs", "Hurls", "Lets it loose",
		"Tosses", "Flings", "Delivers", "Launches", "Sends", "Rifles", "Floats", "Darts",
		"Whips", "Releases", "Unleashes", "Heaves", "Zips", "Bullets", "Spirals", "Wings"}
	if yards > 19 {
		list = append(list, "Chucks", "Bombs", "Cannons", "Rockets", "Launches deep", "Airs it out", "Goes deep with")
	} else if yards < 5 {
		list = append(list, "Flips", "Dumps", "Pitches", "Pops", "Flicks", "Dishes")
	}
	return PickFromStringList(list)
}

func getInterceptVerb() string {
	list := []string{"intercepted", "picked off"}
	return PickFromStringList(list)
}

func getDistance(yards int) string {
	direction := GenerateIntFromRange(1, 3) // 1 == left, 2 == Middle, 3 == right
	dirs := ""
	var dirsList []string
	switch direction {
	case 1:
		dirsList = []string{"towards the left sideline", "towards the left", "to the left side",
			"on the left side of the field", "to the left hash", "toward the left boundary",
			"on the weak side", "to the left edge", "along the left sideline"}
		dirs = "left"
	case 2:
		dirsList = []string{"towards the middle of the field", "towards the middle", "over the middle",
			"up the seam", "between the hash marks", "through the center", "right down the middle",
			"straight ahead", "up the gut", "through the heart of the defense", "down the center",
			"between the numbers", "through the middle corridor", "up the spine"}
	default:
		dirsList = []string{"towards the right sideline", "towards the right", "to the right side",
			"on the right side of the field", "to the right hash", "toward the right boundary",
			"on the strong side", "to the right edge", "along the right sideline"}
		dirs = "right"
	}

	// Add depth-specific descriptions
	if yards > 30 {
		dirsList = append(dirsList, "deep downfield", "way downfield", "far down the field",
			"into the secondary", "deep into coverage", "well beyond the sticks")
		if direction == 1 || direction == 3 {
			dirsList = append(dirsList, "deep to the "+dirs+" side", "far downfield to the "+dirs,
				"deep along the "+dirs+" sideline", "way downfield on the "+dirs)
		} else {
			dirsList = append(dirsList, "deep down the middle", "far up the seam", "deep through the heart")
		}
	} else if yards > 15 {
		dirsList = append(dirsList, "intermediate level", "at the intermediate depth", "to the second level")
		if direction == 1 || direction == 3 {
			dirsList = append(dirsList, "intermediate to the "+dirs, "mid-level on the "+dirs+" side")
		} else {
			dirsList = append(dirsList, "intermediate up the middle", "mid-level through the center")
		}
	} else if yards > 5 {
		dirsList = append(dirsList, "short", "underneath", "at short depth")
	} else {
		dirsList = append(dirsList, "quick", "very short", "just beyond the line", "at the line level")
	}

	return PickFromStringList(dirsList)
}

func CheckSubstring(text, subtext string) bool {
	return strings.Contains(text, subtext)
}

func GetYardsString(yds int8) string {
	yards := " yards. "
	if yds == 1 || yds == -1 {
		yards = " yard. "
	}
	return yards
}
