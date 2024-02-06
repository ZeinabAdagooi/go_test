
SELECT 
    CASE 
        WHEN id % 2 = 1 AND id < (SELECT MAX(id) FROM Seat) THEN 
            (SELECT id FROM Seat WHERE id = t.id + 1)
        WHEN id % 2 = 1 THEN id
        ELSE (SELECT id FROM Seat WHERE id = t.id - 1)
    END AS id , student
FROM Seat t
ORDER BY id;



