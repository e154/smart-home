-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
-- alter table dashboard_tabs
ALTER TABLE dashboard_tabs
    DROP CONSTRAINT if exists dashboard_tabs_background_check,
    ALTER COLUMN background DROP DEFAULT,
    ALTER COLUMN background DROP NOT NULL;
ALTER TABLE dashboard_cards
    DROP CONSTRAINT if exists dashboard_cards_background_check,
    ALTER COLUMN background DROP DEFAULT,
    ALTER COLUMN background DROP NOT NULL;

update dashboard_tabs
set background = null
where background = 'white'
   or background = '';
update dashboard_cards
set background = null
where background = 'white';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


