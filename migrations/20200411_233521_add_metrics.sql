-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table metrics
(
    id          bigserial primary key,
    name        text                     not null,
    description text,
    options     jsonb default '{}'       not null,
    type        text  default 'line'     not null,
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone not null
);

create index name_at_metrics_unq on metrics (name);

create table map_element_metrics
(
    map_element_id bigint null
        CONSTRAINT map_element_id_at_map_element_metric_2_map_elements_fk REFERENCES map_elements (id) ON UPDATE CASCADE ON DELETE CASCADE,
    metric_id      bigint null
        CONSTRAINT metric_id_at_map_element_metric_2_metric_fk REFERENCES metrics (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (map_element_id, metric_id)
);

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
drop table map_element_metrics cascade;
drop table metric_bucket cascade;
drop table metrics cascade;
