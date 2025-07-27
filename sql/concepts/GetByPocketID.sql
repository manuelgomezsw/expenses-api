SELECT c.id,
       c.name,
       c.value,
       c.pocket_id,
       p.name AS 'pocket_name',
       c.payed,
       c.updated_at,
       c.payment_day
FROM concepts c
         JOIN pockets p ON c.pocket_id = p.id
WHERE pocket_id = ?