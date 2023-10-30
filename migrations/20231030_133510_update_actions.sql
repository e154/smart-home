-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table actions
    add area_id     bigint null
        constraint action_2_areas_fk
            references areas
            on update cascade on delete set null,
    add description text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


