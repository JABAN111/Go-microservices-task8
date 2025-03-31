## Просто забавная статистика для просмотра пользы Gin

### Запрос без индекса
```sql
postgres=# EXPLAIN ANALYZE SELECT id 
FROM comics 
WHERE words && ARRAY['window']
;
                                             QUERY PLAN                                             
----------------------------------------------------------------------------------------------------
 Seq Scan on comics  (cost=0.00..234.35 rows=72 width=4) (actual time=0.055..2.262 rows=72 loops=1)
   Filter: (words && '{window}'::text[])
   Rows Removed by Filter: 2996
 Planning Time: 0.181 ms
 Execution Time: 2.282 ms
(5 rows)
```

### c Gin индексом

```sql
postgres=# EXPLAIN ANALYZE SELECT id 
FROM comics 
WHERE words && ARRAY['window']
;
                                                         QUERY PLAN                                                         
----------------------------------------------------------------------------------------------------------------------------
 Bitmap Heap Scan on comics  (cost=13.20..156.01 rows=72 width=4) (actual time=0.049..0.127 rows=72 loops=1)
   Recheck Cond: (words && '{window}'::text[])
   Heap Blocks: exact=59
   ->  Bitmap Index Scan on idx_comics_words  (cost=0.00..13.18 rows=72 width=0) (actual time=0.038..0.039 rows=72 loops=1)
         Index Cond: (words && '{window}'::text[])
 Planning Time: 0.421 ms
 Execution Time: 0.160 ms
(7 rows)
```

