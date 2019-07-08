
-- +migrate Up
alter table articles drop column old_pic_id;

-- +migrate Down
