-- https://stackoverflow.com/questions/8146448/get-the-default-values-of-table-columns-in-postgres
-- https://github.com/volatiletech/sqlboiler/issues/409
--
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
--
-- To mitigate the above situation we fully disallow to set DEFAULT to anything
-- other than the default golang zero value of this type. Otherwise this isssue
-- is fairly hard to diagnose.
--
--
-- If a default value is actually set, we only the respective mapped golang zero value:
-- 0 for all integer types
-- 0.0 for floating point numbers
-- false for booleans
-- "" for strings
-- https://yourbasic.org/golang/default-zero-value/

SELECT
    table_name,
    column_name,
    data_type,
    column_default,
    is_nullable
FROM
    information_schema.columns
WHERE (table_schema) = ('public')
ORDER BY
    table_name,
    column_name;

