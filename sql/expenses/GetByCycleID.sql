SELECT e.id,
       e.name,
       e.value,
       e.cycle_id,
       c.name                            as 'cycle_name',
       p.id                              as 'pocket_id',
       CONCAT(p.name, ' (', c.name, ')') as 'pocket_name',
       e.payment_type_id,
       pt.name                           as 'payment_type_name',
       e.created_at
FROM expenses e
         JOIN cycles c
              ON c.id = e.cycle_id
                  AND c.status = true
         JOIN pockets p
              ON p.id = c.pocket_id
         JOIN payment_types pt
              ON pt.id = e.payment_type_id
WHERE c.id = ?
ORDER BY e.created_at DESC