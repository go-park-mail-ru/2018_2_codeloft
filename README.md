## Многопользовательская браузерная игра Tron Remastered

- [Ссылки](#references)
- [Описание](#description)
- [Участники](#authors)
- [Frontend](#front)
- [API](#API)

## Приложение доступно по хостам: <a name="references"></a>

- [Ветка разработки](https://github.com/frontend-park-mail-ru/2018_2_codeloft/tree/testing)
- [Release ветка](https://github.com/frontend-park-mail-ru/2018_2_codeloft/tree/master)
- [Деплой](https://codeloft-backend.now.sh)

## Описание <a name="description"></a>

2D игра по мотивам легендарной аркады TRON. Игра представляет собой сражение на байках,
оставляющих за собой непреодолимый барьер, столкновение с которым окажется фатальным.
Игроки встречаются на поле в попытках обыграть друг друга не только используя реакцию, но и интеллект.
Хитрая стратегия, опыт, навыки и немного удачи в виде бонусов разбросанных по карте помогут определить кто сильнейший.


## Участники <a name="authors"></a>
- [Анохин Данил](https://github.com/Malefaro), backend
- [Рязанов Максим](https://github.com/RyazMax), backend
- [Дыров Игорь](https://github.com/igor-dyrov), frontend
- [Саркисян Артур](https://github.com/Arthurunique24), fullstack

## Frontend <a name="back"></a>
- [Frontend repository](https://github.com/frontend-park-mail-ru/2018_2_codeloft)

## API
_______
###   **/user**
* **Method GET:** LeaderBoard
** *Take:* get params "page" and "page_size"
** Return: JSON with leaders

* **Method POST:** Registration
** Take: JSON with "login", "password", "email"
** Return: JSON with "user_id","login", "password", "email"

* **Method DELETE:** delete user
** Take: JSON with "login", "password"
** Return: nothing

* **Method PUT:** update user
** Take: JSON with "login", "password" and optional "email", "new_password", "score"
** Return: JSON with "user_id","login", "password", "email" of update user

### **/user/id**
* **Method GET:** get user with id

### **/session**
* **Method GET:** checkAuth with cookie

* **Method POST:** login. set cookie
**Take: JSON with "login", "password"

* **Method DELETE:** logout,delete cookie
**Take: JSON with "login", "password"
