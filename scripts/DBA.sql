/* psql commands */
SELECT pg_size_pretty(pg_relation_size('poster'));
SELECT pg_size_pretty(pg_indexes_size('poster'));
SELECT pg_size_pretty(pg_total_relation_size('poster'));

SELECT
    table_name,
    pg_size_pretty(table_size) AS table_size,
    pg_size_pretty(indexes_size) AS indexes_size,
    pg_size_pretty(total_size) AS total_size
FROM (
    SELECT
        table_name,
        pg_table_size(table_name) AS table_size,
        pg_indexes_size(table_name) AS indexes_size,
        pg_total_relation_size(table_name) AS total_size
    FROM (
        SELECT ('"' || table_schema || '"."' || table_name || '"') AS table_name
        FROM information_schema.tables WHERE table_schema = 'public'
    ) AS all_tables
    ORDER BY total_size DESC
) AS pretty_sizes;

/* Manipultate SQLs */
ALTER TABLE poster RENAME COLUMN created_time TO created_time;

ALTER TABLE poster ADD PRIMARY KEY (email);
ALTER TABLE poster ADD CONSTRAINT username_unique UNIQUE (username);
ALTER TABLE poster DROP CONSTRAINT poster_role_check;
ALTER TABLE poster ALTER COLUMN password TYPE char(60);

DROP INDEX [ CONCURRENTLY] [ IF EXISTS ] username_unique [ CASCADE | RESTRICT ];
ALTER TABLE poster DROP CONSTRAINT username_unique;

SELECT count(*), state FROM pg_stat_activity GROUP BY 2;

/* Operational SQLs */
SELECT * FROM ONLY poster;

SELECT * FROM poster LIMIT 100;
SELECT * FROM poster;
SELECT * FROM poster WHERE email = 'linushung@gmail.com';
SELECT * FROM poster WHERE username = 'linushung';
SELECT * FROM follower WHERE email = 'sabrinaho@gmail.com' AND follower = 'linushung';
-- SELECT pc.*, pi.member1_id, pi.member2_id FROM pair_chatroom pc INNER JOIN pair_info pi on pc.id = pi.room_id WHERE pc.id = '0997693bce954c91a7c3c7971a132bb8_2rm';
-- SELECT room_id, user_id, count(*) FROM chatroom_anonymous_info GROUP BY room_id, user_id HAVING count(*) > 1

/* Ref: RDBMS
1. https://vladmihalcea.com/how-does-a-relational-database-work/
2. https://vladmihalcea.com/a-beginners-guide-to-acid-and-database-transactions/
3. https://tapoueh.org/blog/2018/03/database-normalization-and-primary-keys/
*/

/* Ref: Isolation
1. https://vladmihalcea.com/how-does-mvcc-multi-version-concurrency-control-work/
2. https://vladmihalcea.com/2pl-two-phase-locking/
3. https://vladmihalcea.com/postgresql-triggers-isolation-levels/
4. https://www.postgresql.org/docs/10/transaction-iso.html
*/

/* Discussion: Null with Normalization
1. https://stackoverflow.com/questions/8684001/is-this-a-1nf-failure
2. https://stackoverflow.com/questions/163434/are-nulls-in-a-relational-database-okay
*/
