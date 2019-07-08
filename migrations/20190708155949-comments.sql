-- +migrate Up
-- таблица пользователей
create table if not exists users (
  uuid uuid primary key default uuid_generate_v4 (),
  old_id varchar,
  login varchar not null,
  password varchar not null,
  email varchar not null,
  recovery_hash varchar,
  last_visit timestamp with time zone
);

-- таблица комментов
create table if not exists comments (
  uuid uuid primary key default uuid_generate_v4 (),
  old_id varchar,
  old_post_id varchar,
  old_author_id varchar,
  user_uuid uuid references users (uuid) on delete cascade,
  article_uuid uuid references articles (uuid) on delete cascade,
  message text not null,
  is_visible boolean default true,
  created_at timestamp with time zone default current_timestamp
);

-- создадим индекс на время публикации статьи
create index active_time_indx on articles (active_from);

-- создадим индекс на время публикации комментария
create index created_at_idx on comments (created_at);

-- создадим индекс на логин пользователя
create index user_login_idx on users (login);

-- старый id статьи не обязателен
alter table articles alter column old_id drop not null;

-- В случае добавления пользователя - безусловно хешируем его пароль
-- +migrate StatementBegin
create or replace function hash_password_on_insert ()
  returns trigger
as $$
begin
  new.password = crypt(new.password, gen_salt('bf'));
  return new;
end;
$$
language plpgsql;
-- +migrate StatementEnd

-- Триггер на хеширование при добавлении
create trigger hash_password_on_insert_trigger before
insert
  on users for each row execute procedure hash_password_on_insert ();

-- В случае правки пользователя- хешируем пароль, если он был задан и он не пуст
-- +migrate StatementBegin
create or replace function hash_password_on_update ()
  returns trigger
as $$
begin
  if new.password is not null and new.password != '' and new.password != old.password then
    new.password = crypt(new.password, gen_salt('bf'));
  end if;
  return new;
end;
$$
language plpgsql;
-- +migrate StatementEnd

-- Триггер на хеширование при правке пользователя
create trigger hash_password_on_update_trigger before
update
    on users for each row execute procedure hash_password_on_update ();

-- Функция проверки имеется ли пользователь с таким логином и паролем
-- +migrate StatementBegin
create function user_has_password(varchar, varchar) returns boolean as $$
declare
	has_password boolean;
begin
	select (u.password=crypt($2, u.password))::boolean into has_password from users u where u.login=$1;
	if found
	then
		return has_password;
	else
		return false;
	end if;
end
$$ language 'plpgsql';
-- +migrate StatementEnd

-- +migrate Down
drop trigger if exists hash_password_on_insert_trigger on users cascade;
drop trigger if exists hash_password_on_update_trigger on users cascade;
drop function if exists hash_password_on_insert ();
drop function if exists hash_password_on_update ();
drop function if exists user_has_password ();
drop table if exists users cascade;
drop table if exists comments cascade;

