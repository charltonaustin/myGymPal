INSERT INTO phases (program_id, phase_number)
SELECT p.id, gs.n
FROM programs p
CROSS JOIN generate_series(1, p.num_phases) gs(n)
ON CONFLICT (program_id, phase_number) DO NOTHING;
