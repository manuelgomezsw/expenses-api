SELECT id,
       name,
       status,
       created_at
FROM pockets
WHERE status = true
ORDER BY name