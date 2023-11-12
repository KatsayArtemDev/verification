create database verification;

create extension if not exists "uuid-ossp";

drop table pins
;

set timezone to 'UTC';

create table pins(
  user_id uuid not null primary key,
  sent_at timestamp,
  pin text not null
)
;

create table attempts(
  user_id uuid not null,
  attempts integer default 1,
  primary key (user_id)
)
;

create table blocks(
  user_id uuid not null,
  blocked_at timestamp default (current_timestamp at time zone 'UTC'),
  primary key (user_id)
)
;
