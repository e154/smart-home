-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table metrics
(
    id            bigserial primary key,
    map_device_id bigint                   null
        CONSTRAINT map_device_at_metrics_2_map_devices_fk REFERENCES map_devices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    name          text                     not null,
    description   text,
    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone not null
);

create index name_at_metrics_unq on metrics (name);

create table metric_bucket
(
    value     jsonb                    not null,
    metric_id bigint                   null
        CONSTRAINT metric_id_at_metric_bucket_2_metric_fk REFERENCES metrics (id) ON UPDATE CASCADE ON DELETE CASCADE,
    time      timestamp with time zone not null,
    PRIMARY KEY (time, metric_id)
);

create index time_at_metric_bucket_idx on metric_bucket (time);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table metric_bucket cascade;
drop table metrics cascade;
