-- +migrate Up
alter table users drop column old_id;
alter table users add constraint email_is_uniq unique (email);
alter table comments drop column old_id;
alter table comments drop column old_post_id;
alter table comments drop column old_author_id;
alter table comments alter column user_uuid set not null;
alter table comments alter column article_uuid set not null;
alter table articles drop column old_id;
alter table articles add column image varchar;

-- +migrate Down
