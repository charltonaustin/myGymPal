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
