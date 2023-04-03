ALTER TABLE IF EXISTS {{ params.target_schema }}.{{ params.tablename }}
    RENAME TO {{ params.tablename }}_old;

ALTER INDEX IF EXISTS {{ params.target_schema }}.{{ params.tablename }}_pkey
    RENAME TO {{ params.tablename }}_old_pkey;

{% if params.index_names is defined %}
{% for index_name in params.index_names %}
ALTER INDEX IF EXISTS {{ params.target_schema }}.{{ index_name }}
    RENAME TO {{ index_name }}_old;
{% endfor %}
{% endif %}

ALTER TABLE {{ params.staging_schema }}.{{ params.tablename }}
    SET SCHEMA {{ params.target_schema }};

DROP TABLE IF EXISTS {{ params.target_schema }}.{{ params.tablename }}_old;

