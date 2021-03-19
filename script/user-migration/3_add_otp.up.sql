ALTER TABLE users ADD COLUMN use_otp boolean null;

CREATE TABLE otps (
  id          serial primary key,
  secret      varchar(250) not null,
  expiration  timestamp not null,
  created_at  timestamp not null,
  updated_at  timestamp not null,
  deleted_at  timestamp null,
  username    varchar(255) not null,
  foreign key (username) references users (username)
);

