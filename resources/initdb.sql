CREATE EXTENSION  IF NOT EXISTS citext;

create table if not exists users (
  id bigserial not null primary key,
  login citext unique ,
  password varchar(30),
  email varchar(30),
  score bigint default 0,
  lang varchar(2) default 'en'
);

create table if not exists sessions (
  value text unique,
  id int references users(id) on delete cascade
);

insert into users(login,password,email) values ('kek','qwerty12345','kek@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek2','qwerty12345','kek2@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek3','qwerty12345','kek3@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek4','qwerty12345','kek4@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek5','qwerty12345','kek5@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek6','qwerty12345','kek6@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek7','qwerty12345','kek7@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek8','qwerty12345','kek8@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek9','qwerty12345','kek9@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek10','qwerty12345','kek10@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek11','qwerty12345','kek11@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek12','qwerty12345','kek12@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek13','qwerty12345','kek13@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek14','qwerty12345','kek14@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek15','qwerty12345','kek15@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek16','qwerty12345','kek16@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek17','qwerty12345','kek17@mail.ru') on CONFLICT do nothing;
insert into users(login,password,email) values ('kek18','qwerty12345','kek18@mail.ru') on CONFLICT do nothing;



update users set score=20 where id = 1;
update users set score=15 where id = 2;
update users set score=30 where id = 3;
update users set score=110 where id = 6;

insert into sessions(value, id) values ('asdsa', 3) on CONFLICT do nothing;
insert into sessions(value, id) values ('asdsa2', 4) on CONFLICT do nothing;
insert into sessions(value, id) values ('asdsa3', 6) on CONFLICT do nothing;

 -- psql -U kexibq -d codeloft -a -f initdb.sql

