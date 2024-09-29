create table if not exists device_settings (
  device_id bigint primary key,
  heating_status boolean not null default false,
  updated_at timestamp with time zone not null default now()
);