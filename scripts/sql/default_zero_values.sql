-- Errors if DEFAULT values for certain column data_types
-- is NOT set to the golang zero value
--
-- https://github.com/volatiletech/sqlboiler/issues/409
-- https://github.com/volatiletech/sqlboiler/issues/237
-- https://golang.org/ref/spec#The_zero_value
--
-- https://github.com/volatiletech/sqlboiler#diagnosing-problems
-- A field not being inserted (usually a default true boolean), boil.Infer
-- looks at the zero value of your Go type (it doesn’t care what the default
-- value in the database is) to determine if it should insert your field or
-- not. In the case of a default true boolean value, when you want to set it to
-- false; you set that in the struct but that’s the zero value for the bool
-- field in Go so sqlboiler assumes you do not want to insert that field and
-- you want the default value from the database. Use a whitelist/greylist to
-- add that field to the list of fields to insert.
--
-- boil.Infer() assumes all SQL defaults are Go zero value
--
-- To mitigate the above situation we disallow setting DEFAULT to anything
-- other than the default golang zero value of this type. Otherwise this issue
-- is fairly hard to debug (boil.Infer() does not insert DEFAULT as expected).
--
-- If a default value is actually set, we only allow the respective mapped golang zero value:
-- 0 for all integer types
-- 0.0 for floating point numbers
-- false for booleans
-- "" for strings
-- nil for pointers
--
-- https://stackoverflow.com/questions/8146448/get-the-default-values-of-table-columns-in-postgres
-- https://dba.stackexchange.com/questions/205471/why-does-information-schema-have-yes-and-no-character-strings-rather-than-bo
CREATE OR REPLACE FUNCTION check_default_go_sql_zero_values ()
    RETURNS SETOF information_schema.columns
    AS $BODY$
BEGIN
    RETURN QUERY
    SELECT
        *
    FROM
        information_schema.columns
    WHERE (table_schema = 'public'
        AND column_default IS NOT NULL)
        AND is_nullable = 'NO'
        AND ((data_type = 'boolean'
                AND column_default <> 'false')
            OR (data_type IN ('char', 'character', 'varchar', 'character varying', 'text')
                AND column_default NOT LIKE concat('''''', '::%'))
            OR (data_type IN ('smallint', 'integer', 'bigint', 'smallserial', 'serial', 'bigserial')
                AND (column_default <> '0'
                    AND column_default NOT LIKE 'nextval(%'))
            OR (data_type IN ('decimal', 'numeric', 'real', 'double precision')
                AND column_default <> '0.0'));
END
$BODY$
LANGUAGE plpgsql
SECURITY DEFINER;

CREATE OR REPLACE FUNCTION default_zero_values ()
    RETURNS void
    AS $$
DECLARE
    item record;
BEGIN
    FOR item IN
    SELECT
        table_name,
        column_name,
        data_type,
        column_default,
        is_nullable
    FROM
        check_default_go_sql_zero_values ()
        LOOP
            RAISE WARNING ' %.% % : INVALID DEFAULT ''%''', item.table_name, item.column_name, item.data_type, item.column_default USING HINT = to_json(item);
END LOOP;
    IF FOUND THEN
        RAISE EXCEPTION 'NOT NULL columns require the respective go zero value () AS their DEFAULT value or no DEFAULT at all'
            USING HINT = '0 for integer types, 0.0 for floating point numbers, false for booleans, "" for strings';
        END IF;
END;
$$
LANGUAGE plpgsql;

SELECT
    *
FROM
    default_zero_values ();

