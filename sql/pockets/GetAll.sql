SELECT p.id,
       p.name,
       p.status,
       IFNULL(SUM(c.value), 0) AS TOTAL_AMOUNT,
       p.created_at
FROM pockets p
         LEFT JOIN concepts c
                   ON p.id = c.pocket_id
WHERE p.status = true
GROUP BY
    p.id,
    p.name,
    p.status,
    p.created_at
ORDER BY TOTAL_AMOUNT DESC