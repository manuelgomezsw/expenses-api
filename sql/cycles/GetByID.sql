SELECT c.id,
       c.pocket_id,
       p.name,
       c.name,
       c.budget,
       c.date_init,
       c.date_end,
       c.status,
       c.created_at
FROM cycles c
         JOIN pockets p ON c.pocket_id = p.id
WHERE c.id = ?