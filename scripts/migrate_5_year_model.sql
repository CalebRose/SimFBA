-- =============================================================================
-- Migration Script: NCAA 5-Year Eligibility Model
-- Run this script against the SimFBA database after the active season concludes.
-- =============================================================================

-- 1. Clear redshirt statuses for all active college players
-- Under the 5-year model, players' numerical 'year' (1 to 5) maps directly to:
-- Year 1: Freshman (Fr)
-- Year 2: Sophomore (So)
-- Year 3: Junior (Jr)
-- Year 4: Senior (Sr)
-- Year 5: 5th Year (5th)
UPDATE college_players 
SET is_redshirt = FALSE, is_redshirting = FALSE 
WHERE id > 0;

-- 2. Remove all 'No Redshirt' transfer portal promises
DELETE FROM college_promises 
WHERE promise_type = 'No Redshirt';

-- 3. Verification query: Confirm no active players have invalid years
SELECT id, player_id, first_name, last_name, team_abbr, year, is_redshirt, is_redshirting
FROM college_players
WHERE year < 1 OR year > 5;

