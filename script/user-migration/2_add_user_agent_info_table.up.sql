CREATE TABLE agent_infos (
  id          serial primary key,
  address     varchar (50) not null,
  client      varchar(255) not null,
  created_at  timestamp not null,
  updated_at  timestamp not null,
  deleted_at  timestamp null,
  user_id     integer,
  foreign key (user_id) references users (id)
);