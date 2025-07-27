SELECT p.id,
       p.name,
       p.status,
       SUM(c.value) AS TOTAL_AMOUNT,
       p.created_at
FROM pockets p
         INNER JOIN concepts c
                    ON p.id = c.pocket_id
WHERE p.id = ?
GROUP BY
    p.id,
    p.name,
    p.status,
    p.created_at
ORDER BY name