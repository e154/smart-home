-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table alexa_skills
    drop column on_launch cascade;
alter table alexa_skills
    drop column on_session_end cascade;
ALTER TABLE alexa_skills
    ADD COLUMN script_id bigint default null;
ALTER TABLE alexa_skills
    ADD CONSTRAINT alexa_skills_2_scripts_fk
        FOREIGN KEY (script_id) REFERENCES scripts (id) ON DELETE SET NULL ON UPDATE CASCADE;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


