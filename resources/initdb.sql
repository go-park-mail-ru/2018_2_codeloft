create table if not exists users (
  id bigserial not null primary key,
  login varchar(30),
  password varchar(30),
  email varchar(30),
  score int
);

insert into users(login,password,email,score) values ('kek','qwerty12345','kek@mail.ru',0);
insert into users(login,password,email,score) values ('kek2','qwerty12345','kek@mail.ru',0);
insert into users(login,password,email,score) values ('kek3','qwerty12345','kek@mail.ru',0);
insert into users(login,password,email,score) values ('kek4','qwerty12345','kek@mail.ru',0);

