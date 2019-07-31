-- +migrate Up
-- переводим все емейлы в lowercase
update users u set email = lower(u.email);

-- емейл и логин уникальны
alter table users add constraint uniq_login unique (login);
alter table users add constraint uniq_email unique (email);

-- добавим индекс по емейлу, ибо по нему будет авторизация
create index if not exists email_idx on users (email);

-- при добавлении пользователя емейл приводится к нижнему регистру

-- +migrate StatementBegin
create or replace function email_to_lower() returns trigger as $body$
begin
  NEW.email = lower( NEW.email );
  return NEW;
end;
$body$ language plpgsql;
-- +migrate StatementEnd

create trigger email_to_lower before insert or update on users
  for each row execute procedure email_to_lower ();

-- функция для аутентификации по Email
-- +migrate StatementBegin
create or replace function email_has_password(eml varchar, pass varchar) returns boolean as $body$
declare 
  has_password boolean;
  pass_hash varchar;
begin
  select u.password into pass_hash from users u where u.email = eml;
  if not found then
    raise exception 'user_not_found_by_email' using errcode = 'no_data_found';
  end if;
  if pass_hash = crypt(pass, pass_hash) then
    return true;
  end if;
  return false;
end;
$body$ language plpgsql;
-- +migrate StatementEnd

-- таблица токенов
create table if not exists tokens (
  user_uuid uuid not null references users(uuid) on delete cascade,
  access_token uuid not null default uuid_generate_v4 (),
  refresh_token uuid not null default uuid_generate_v4 (),
  updated_at timestamp with time zone default current_timestamp
);
create index if not exists access_token_index on tokens(access_token);
create index if not exists refresh_token_index on tokens(refresh_token);
create trigger update_updated_at before update
  on tokens for each row execute procedure update_updated_at_column ();

-- функция проверки действительности access_token
-- возвращает true если токен актуален, false если устарел, exception если не найден
-- +migrate StatementBegin
create or replace function is_valid_access_token(token uuid) returns boolean as $body$
declare
  upd timestamp with time zone;
begin
  select updated_at into upd from tokens t where t.access_token = token;
  if not found then
    return false;
  end if;
  if upd + interval '5 minutes' > current_timestamp then
    return true;
  end if;
  return false;
end;
$body$ language plpgsql;
-- +migrate StatementEnd

--функция обновления токена по refresh_token
-- +migrate StatementBegin
create or replace function refresh_tokens(rtoken uuid, out a_t uuid, out r_t uuid)as $body$
declare
  tkns record;
begin
  delete from tokens t where t.refresh_token = refresh_token and updated_at + interval '6 months' < current_timestamp;
  begin
    update tokens t
      set 
        refresh_token = uuid_generate_v4 (),
        access_token = uuid_generate_v4 ()
      where
        t.refresh_token = rtoken
      returning 
        access_token, refresh_token
    into strict tkns;
  exception
    when no_data_found then raise exception 'token_not_found' using errcode = 'no_data_found';
  end;
  a_t := tkns.access_token;
  r_t := tkns.refresh_token;
end;
$body$ language plpgsql;
-- +migrate StatementEnd

-- функция создания новой пары токенов
-- +migrate StatementBegin
create or replace function new_tokens(userid uuid, out a_t uuid, out r_t uuid) as $body$
declare
  tkns record;
  user_exists boolean; 
begin
  begin
    insert into tokens(user_uuid)
    values (userid)
    returning access_token, refresh_token into strict tkns;
  exception
    when no_data_found then raise exception 'user_not_found';
    when foreign_key_violation then raise exception 'user_not_found';
  end;
  a_t := tkns.access_token;
  r_t := tkns.refresh_token;
end;
$body$ language plpgsql;
-- +migrate StatementEnd

-- функция логаута для отдельного токена
-- +migrate StatementBegin
create or replace function logout_token(token uuid) returns void as $body$
begin
  delete from tokens t where t.access_token = token;
  return;
end;
$body$ language plpgsql;
-- +migrate StatementEnd

-- функция для логаута всех токенов пользователя
-- +migrate StatementBegin
create or replace function logout_user(uid uuid) returns void as $body$
begin
  delete from tokens t where t.user_uuid = uid;
  return;
end;
$body$ language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
