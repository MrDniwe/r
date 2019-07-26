-- +migrate Up
-- хеш восстановления это теперь uuid
alter table users drop column recovery_hash;
alter table users add column recovery_hash uuid;
-- записываем время генерации хеша
alter table users add column hash_generated_at timestamp with time zone;

-- +migrate StatementBegin
create or replace function create_recovery_hash(eml varchar) returns uuid as $body$
declare
  user_uuid uuid;
  can_gen_code boolean;
  hash uuid;
begin
  -- проверяем есть ли пользователь с такой почтой
  select uuid into user_uuid from users u where lower(eml) = lower(u.email);
  if user_uuid is null then
    raise exception 'user_not_found_by_email' using errcode = 'no_data_found';
  end if;
  -- проверяем как давно ему генерился код восстановления
  select hash_generated_at < current_timestamp - interval '1 minute' into can_gen_code
    from users u
    where u.uuid = user_uuid;
  if can_gen_code = false then
    raise exception 'codegen_request_to_soon' using errcode = 'data_exception';
  end if;
  -- генерим новый код и возвращаем его
  update users u
  set
    recovery_hash = uuid_generate_v4(),
    hash_generated_at = current_timestamp
  where u.uuid = user_uuid
  returning recovery_hash into hash;
  return hash;
end;
$body$ language plpgsql;
-- +migrate StatementEnd


-- +migrate Down
