-- +migrate Up
-- модули для работы с uuid и криптографией
create extension if not exists pgcrypto;
create extension if not exists "uuid-ossp";

-- таблица для хранения статей
create table if not exists articles (
  uuid uuid primary key default uuid_generate_v4 (),
  old_id varchar not null,
  old_pic_id varchar,
  title varchar not null,
  lead text,
  body text,
  active_from timestamp with time zone,
  views int,
  is_visible boolean default false,
  created_at timestamp with time zone not null default current_timestamp,
  updated_at timestamp with time zone not null default current_timestamp
);

-- Общая для всех триггерная функция, обновляющая дату изменения
-- +migrate StatementBegin
create or replace function update_updated_at_column ()
    returns trigger
as $$
begin
    NEW.updated_at = current_timestamp;
    return NEW;
end;
$$
language plpgsql;
-- +migrate StatementEnd

-- Триггер для сохранения времени изменения роли
create trigger update_updated_at before
update
    on articles for each row execute procedure update_updated_at_column ();

-- +migrate Down
drop trigger if exists update_updated_at on articles cascade;
drop function if exists update_updated_at_column ();
drop table articles if exists cascade;
