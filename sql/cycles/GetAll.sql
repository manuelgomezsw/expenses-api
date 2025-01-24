SELECT c.id,
       p.name as 'pocket_name',
       c.name as 'cycle_name',
       c.budget,
       c.date_init,
       c.date_end,
       c.status,
       c.created_at
FROM cycles c
JOIN pockets p ON c.pocket_id = p.id
ORDER BY c.budget DESC