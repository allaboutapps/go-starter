-- Errors if FKs do not have an index
CREATE OR REPLACE FUNCTION fk_missing_index ()
    RETURNS void
    AS $$
DECLARE
    item record;
BEGIN
    FOR item IN
    SELECT
        c.conrelid::regclass AS "table",
        /* list of key column names in order */
        string_agg(a.attname, ',' ORDER BY x.n) AS columns,
        pg_catalog.pg_size_pretty(pg_catalog.pg_relation_size(c.conrelid)) AS size,
        c.conname AS constraint,
        c.confrelid::regclass AS referenced_table
    FROM
        pg_catalog.pg_constraint c
        /* enumerated key column numbers per foreign key */
    CROSS JOIN LATERAL unnest(c.conkey)
    WITH ORDINALITY AS x (attnum, n)
    /* name for each key column */
    JOIN pg_catalog.pg_attribute a ON a.attnum = x.attnum
        AND a.attrelid = c.conrelid
WHERE
    NOT EXISTS
    /* is there a matching index for the constraint? */
    (
        SELECT
            1
        FROM
            pg_catalog.pg_index i
        WHERE
            i.indrelid = c.conrelid
            /* the first index columns must be the same as the
             key columns, but order doesn't matter */
            AND (i.indkey::smallint[])[0:cardinality(c.conkey) - 1] @> c.conkey)
        AND c.contype = 'f'
    GROUP BY
        c.conrelid,
        c.conname,
        c.confrelid
    ORDER BY
        pg_catalog.pg_relation_size(c.conrelid) DESC LOOP
            RAISE WARNING 'CREATE INDEX "idx_%_fk_%" ON "%" ("%");', item.table, item.columns, item.table, item.columns USING HINT = to_json(item);
        END LOOP;
    IF FOUND THEN
        RAISE EXCEPTION ' We require ALL FOREIGN keys TO have an INDEX defined. ';
    END IF;
END;
$$
LANGUAGE plpgsql;

SELECT
    *
FROM
    fk_missing_index ();

