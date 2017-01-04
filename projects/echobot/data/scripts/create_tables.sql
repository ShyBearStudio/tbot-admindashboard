drop table if exists echobot_chats cascade;

create table echobot_chats (
  id         integer unique,
  type       varchar(32) not null,
  title      varchar(255),
  username   varchar(255),
  firstname  varchar(255),
  lastname   varchar(255),
  created_at timestamp not null,
  active     boolean not null);
