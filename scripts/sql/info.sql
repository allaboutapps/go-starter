-- This file just queries for the currently applied tables/columns as an overview
SELECT
    table_name,
    column_name,
    data_type,
    column_default,
    is_nullable
FROM
    information_schema.columns
WHERE (table_schema = 'public')
ORDER BY
    table_name,
    is_nullable,
    column_name;

