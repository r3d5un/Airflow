SELECT *
FROM
    {{ params.schema }}.{{ params.table }}
LIMIT 10;
