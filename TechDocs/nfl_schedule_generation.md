# Generating NFL schedules for the simulation.

Title is straightforward. We want to create a function that will create a full season's worth of games for SimNFL based on the NFL's schedule template format. Each team needs 17 games within an 18 week span.

## Scheduling Formula

The scheduling formula should be based on the NFL's current scheduling formula as detailed in this [wikipedia article](https://en.wikipedia.org/wiki/NFL_regular_season#Current_scheduling_formula), which is also explained below. It is a 17-game regular season using a predetermined formula.

Each team plays twice against each of the other three teams in its division: once at home, and once on the road (six games).
Each team plays once against each of the four teams from a predetermined division (based on a three-year rotation) within its own conference: two at home, and two on the road (four games).
Each team plays once against one team from the remaining two divisions within its conference that finished in the same placement in the final divisional standings in the prior season:[b] one at home, one on the road (two games).
Each team plays once against each of the four teams from a predetermined division (based on a four-year rotation) in the other conference: two at home, and two on the road (four games).
Each team also plays an extra interconference "17th game" against one team from the division in the other conference it played two years before that. Additionally, the opponent is decided based on where the two teams finish in their respective divisions in the previous season[c] (one game).

Under this formula, all teams are guaranteed to play every other team in their own conference at least once every three years, and to play every team in the other conference at least once every four years. The formula also guarantees a similar schedule for every team in a division each season, as all four teams will play fourteen out of their seventeen games against common opponents or each other.

# Game Times

The Majority of NFL games will likely take place on the 'Sunday Noon' and 'Sunday Afternoon' timeslots in the simulation. There are a few exceptions for this rule:

- One 'Sunday Night Football' game each week
- One 'Monday Night Football' Game each week
- One 'Thursday Night Football' game each week.

## Thanksgiving Day games

During week 13 of the simulation, three Thursday Night Football games will be ran & simulated that week. This is the only exception to the game times rules above. The Thanksgiving Day matchups should be division rivalry games. Here are the strict requirements:

- The Dallas Cowboys need to play a division rival during week 13 on Thursday Night Football. Can be either Home or Away.
- The Detroit Lions need to play a division rival during week 13 on Thursday Night Football. Can be either Home or Away.
- Any other division rivalry matchup can be scheduled for Thursday Night Football during Week 13. Examples of potential rivalry matchups below:
  -- Baltimore Ravens vs Pittsburgh Steelers
  -- New England Patriots vs Buffalo Bills
  -- Kansas City Chiefs vs Denver Broncos
  -- Seattle Seahawks vs San Francisco 49ers
  -- Los Angeles Rams vs San Francisco 49ers
  -- Los Angeles Chargers vs Las Vegas Raiders
  -- Miami Dolphins vs New York Jets
  -- New York Jets vs New England Patriots
  -- Atlanta Falcons vs Carolina Panthers
  -- New Orleans Saints vs Atlanta Falcons
  -- Tampa Bay Buccaneers vs New Orleans Saints

# NFLGame struct requirements

We need the following requirements for each generated record

- WeekID: dataspecific week ID for the season. Take the timestamps' current seasonID and conduct the following formula: (seasonID+2020-2000)\*100 + week

Example: (7 + 2020) = 2027 - 2000 = 27 \* 100 = 2700 + 1 = 2701. That would be the weekID for week 1 of the 2027 season, season ID being 7.

- Week: The week of the season in place. Should be 1-22. Regular season is 1-18.

- HomeTeamID: The Home Team's primary key ID.
- HomeTeam: Can be the abbreviation or TeamName of the home team
- AwayTeamID: The Away Team's primary key ID
- AwayTeam: can be the abbreviation or TeamName of the away team
- HomeTeamCoach: If the home team has a NFLCoachName that is not "" or "AI", place the value here. Otherwise use the NFLOwnerName value as HomeTeamCoach.
- AwayTeamCoach: If the away team has a NFLCoachName that is not "" or "AI", place the value here. Otherwise use the NFLOwnerName value as AwayTeamCoach.
- TimeSlot. If the team is based in the Western or Midwestern US (past the Mississippi River) United States, Sunday Afternoon Timeslot. Otherwise, Sunday Noon Timeslot. Randomly select a game in each weekly slate as a "Thursday Night Football" timeslot, "Sunday Night Football" timeslot, and "Monday Night Football" timeslot
- StadiumID: the StadiumID of the home team's arena. Look at the Stadium.go struct for more info. Use a map reference for each stadium's home team ID if that helps (NFL preferably.)
- Stadium: stadium name
- City: Home team's city
- State: Home team's state
- Weather related data: Look at WeatherManager.go for how it handles eather data. Don't worry about weather related data until all of the game records are created.
- Is Conference: If both matching teams are in the same conference (ConferenceID), set to true.
- IsDivisional: If both teams are part of the same division (DivisionID), set to true.
- Preseason Game: We should setup a structure to allow teams to schedule preseason games but I'll setup a different tech doc for that later.
- Home Previous Bye: If the home team has a bye week the previous week of this matchup, set to true.
- Away Previous Bye: If the away team has a bye week the previous week of this matchup, set to true.

# Notes

- Use the previous season ID for when getting all NFL standings for generating the current season schedule.
