# My Gym Pal — User Stories

## Assumptions & Scope
- Multi-user: anyone can create an account and use the app
- Open source project
- Mobile-first: users log sets in the gym on their phone
- Training programs group phases and weeks; each phase is 8 weeks, ~6.5 phases per year (this can be user defined)
- Each phase has its own rep range (e.g. 14-12 Phase 2, 10–12 in Phase 3, 8–10 in Phase 4); sessions use the rep range of their phase
- The last week of every phase is a deload week (same exercises and weight, reps reduced by 2 below the phase minimum) -- the number of weeks for  phase can be user defined
- Users can create reusable workout templates; templates are shared across all users
- Weight units are lb or kg, set per user
- Goal weight carries over between phases unchanged; the user adjusts based on performance
- Exercises are either weighted (goal weight + reps) or bodyweight (goal reps only); abs follow the same distinction
- Weight adjustment logic applies only to weighted exercises
- Weighted exercises auto-adjust goal weight; bodyweight exercises auto-adjust goal reps — both use the same trigger rules
- Progression increases after 3 consecutive non-deload sessions hitting the max; decreases after failing the minimum once in any non-deload session
- Deload sessions (the last week in a phase) are excluded from all weight adjustment logic
- Both adjustments are applied when the next session is created

---

## Authentication

~~**US-01** — As a new user, I want to create an account with a username and password so that my workout data is private and tied to me.~~

~~**US-02** — As a returning user, I want to log in so that I can access my workout history and templates.~~

~~**US-03** — As a logged-in user, I want to log out so that my account is secure on shared devices.~~

---

## Account Settings

~~**US-04** — As a user, I want to choose my preferred weight unit (lb or kg) in my account settings so that all weights are displayed and entered in the unit I use.~~

---

## Training Programs

~~**US-05** — As a user, I want to create a training program with a name, start date, and number of phases so that my workouts are organized into a structured plan.~~

~~**US-06** — As a user, I want each phase in my program to contain a user defined amount of weeks with a default of 8 so that my training follows a consistent cycle.~~

~~**US-07** — As a user, I want to define a rep range (min and max) per phase (e.g. 10–12 for Phase 3, 8–10 for Phase 4) so that the target reps automatically adjust as I move through the program.~~

**US-08** — As a user, I want week 8 of each phase to be automatically designated as a deload week so that recovery is built into my program.

**US-09** — As a user, I want deload week sessions to automatically set the target reps to 2 below the phase minimum (e.g. 8 reps when the phase range is 10–12) at the current goal weight so that I don't have to manually adjust anything for deload week.

**US-10** — As a user, I want to view the current phase and week I am in within a program so that I always know where I am in my training cycle.

**US-11** — As a user, I want to see an overview of all phases and weeks in my program so that I can track my progress through the full plan.

---

## Workout Templates

~~**US-12** — As a user, I want to create a workout template with a name, focus (e.g. "Upper Body A"), and a list of exercises (with target weights and rep ranges) so that I can reuse it without re-entering everything each session.~~

**US-13** — As a user, I want to define abs, cardio, and stretch blocks in my template so that the full structure of my session is captured.

**US-14** — As a user, I want to browse templates created by other users so that I can discover new workout structures without building from scratch.

**US-15** — As a user, I want to copy another user's template into my own library so that I can use it as-is or customize it for my needs.

**US-16** — As a user, I want to edit an existing template so that I can update it as my program evolves.

**US-17** — As a user, I want to delete a template I no longer use so that my template list stays clean.

---

## Workout Sessions

~~**US-18** — As a user, I want to start a new workout session within my training program so that the correct phase, week, and workout number are set automatically.~~

~~**US-19** — As a user, I want the session to be pre-filled from the associated template so that I only need to log what I actually did.~~

**US-20** — As a user, I want the session form to be easy to use on a phone so that I can log sets quickly between exercises at the gym.

---

## Exercises & Sets

**US-21** — As a user, I want to designate each exercise in a template as either weighted or bodyweight so that the logging form shows the right fields for each type.

**US-22** — As a user, I want to log a weighted set with the weight used and reps completed so that I have a full record of what I actually did.

**US-23** — As a user, I want to log a bodyweight set with only the reps completed so that I am not required to enter a weight that doesn't apply.

**US-24** — As a user, I want to see my goal weight and rep range (for weighted) or goal reps (for bodyweight) displayed next to each set as I log so that I can stay on target without referencing old notes.

**US-25** — As a user, I want to record a weight drop within a single weighted set (e.g. 55 lb x 4 reps → 45 lb x 6 reps) so that I can accurately capture when I can't maintain the target weight throughout a set.

---

## Abs

**US-26** — As a user, I want to log a bodyweight abs exercise with a goal rep count and record actual reps per set so that my ab work is tracked alongside the rest of the session.

**US-27** — As a user, I want to log a weighted abs exercise with a goal weight and rep count and record the actual weight and reps per set so that weighted ab work is tracked the same way as other weighted exercises.

---

## Cardio

**US-28** — As a user, I want to log a cardio activity with a type (e.g. fartlek, steady state), a goal duration, and the actual duration completed so that I can track both intent and reality.

**US-29** — As a user, I want to record a partial cardio session (e.g. 25 of 30 mins) without it being treated as an error so that incomplete efforts are captured 
---

## Stretching

**US-30** — As a user, I want to log a stretch block with a goal hold time per muscle so that my cooldown is tracked as part of the full session.

---

## Progression (Weighted & Bodyweight)

**US-31** — As a user, when I complete all sets at the maximum rep of my phase range for 3 consecutive non-deload sessions on a weighted exercise, I want the app to automatically increase the goal weight when my next session is created so that progression happens without any manual effort.

**US-32** — As a user, when I fail to hit the minimum reps on any set of a weighted exercise in a non-deload session (even once), I want the app to automatically decrease the goal weight when my next session is created so that I am not stuck failing the same weight repeatedly.

**US-33** — As a user, I want deload sessions to be excluded from all weight adjustment logic so that intentionally reduced deload reps do not incorrectly trigger an increase or decrease.

**US-34** — As a user, I want the default weight adjustment increment (both up and down) to be 5 lbs so that I don't have to configure anything to get started.

**US-35** — As a user, I want to set a custom weight adjustment increment per exercise (e.g. 2.5 lbs for isolation movements, 10 lbs for compound lifts) so that adjustments are appropriate for each movement.

**US-36** — As a user, I want to be shown a clear indicator on the completed session when a weight change was triggered (increase or decrease) so that I know what changed and what the new goal weight is for next time.

**US-37** — As a user, I want the goal weight to carry over unchanged when I move to a new phase so that I continue from where I left off and let performance drive any adjustments.

**US-38** — As a user, when I complete all sets at the maximum rep of my phase range for 3 consecutive non-deload sessions on a bodyweight exercise, I want the app to automatically increase the goal reps when my next session is created so that bodyweight exercises progress the same way as weighted ones.

**US-39** — As a user, when I fail to hit the minimum reps on any set of a bodyweight exercise in a non-deload session (even once), I want the app to automatically decrease the goal reps when my next session is created so that my targets stay achievable.

**US-40** — As a user, I want the default rep adjustment increment (both up and down) for bodyweight exercises to be 1 rep so that progression is gradual by default.

**US-41** — As a user, I want to set a custom rep adjustment increment per bodyweight exercise so that I can progress faster or slower depending on the movement.

**US-42** — As a user, I want to be shown a clear indicator on the completed session when a rep goal change was triggered (increase or decrease) on a bodyweight exercise so that I know what changed and what the new goal is for next time.

---

## History & Progress

**US-43** — As a user, I want to view a list of my past workout sessions so that I can review what I've done.

**US-44** — As a user, I want to view the details of a past session so that I can see exactly what weight and reps I did.

**US-45** — As a user, I want to see how my weight/reps for a given exercise have trended over time so that I can measure progress.
