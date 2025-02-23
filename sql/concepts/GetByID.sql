SELECT c.id,
       c.name,
       c.value,
       c.pocket_id,
       p.name AS 'pocket_name',
       c.payed,
       c.updated_at
FROM concepts c
         JOIN pockets p ON c.pocket_id = p.id
WHERE c.id = ?