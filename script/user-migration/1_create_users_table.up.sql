CREATE TABLE users (
  id          serial primary key,
  status      smallint default 0,
  username    varchar(255) not null UNIQUE,
  password    varchar(60) not null,
  created_at  timestamp not null,
  updated_at  timestamp not null,
  deleted_at  timestamp null
);

  CREATE TABLE action_links (
  id          serial primary key,
  status      smallint not null,
  token       varchar(40) not null,
  type        smallint not null,
  expiration  timestamp not null,
  created_at  timestamp not null,
  updated_at  timestamp not null,
  deleted_at  timestamp null,
  user_id     integer,
  foreign key (user_id) references users (id)
);