SELECT id,
       name,
       value,
       pocket_id,
       payed,
       updated_at
FROM concepts
WHERE pocket_id = ?