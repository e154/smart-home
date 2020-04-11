-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type tsdb_type as enum ('counter', 'gauge', 'histogram', 'summaries');

create table ts_metric
(
    id            BIGSERIAL PRIMARY KEY,
    map_device_id bigint                   null
        CONSTRAINT map_device_at_tsdb_2_map_devices_fk REFERENCES map_devices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    type          tsdb_type                not null,
    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone not null
);

create table ts_bucket
(
    value   jsonb                    not null,
    tsdb_id bigint                   null
        CONSTRAINT tsdb_id_tsdb_bucket_2_tsdb_fk REFERENCES ts_metric (id) ON UPDATE CASCADE ON DELETE CASCADE,
    time    timestamp with time zone not null
);

create index time_at_tsdb_bucket_idx on ts_bucket (time);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table ts_bucket cascade;
drop table ts_metric cascade;

