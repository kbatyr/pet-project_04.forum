# Forum by kbatyr
*******************
This project consists in creating a web forum that allows :

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

## Usage

To run project: `make run`

**Additional commands:**

* Build project: `make build`
* Run in docker container: `make docker`
* Stop docker: `make stop`


#### SQLite

In order to store the data in forum (like users, posts, comments, etc.) I used the database library SQLite.

SQLite is a popular choice as an embedded database software for local/client storage in application software such as web browsers. It enables to create a database as well as controlling it by using queries.

#### Authentication

In this segment the client can `register` as a new user on the forum, by inputting their credentials. Also have to create a `login session` to access the forum and be able to add posts and comments.

I used cookies to allow each user to have only one opened session. Each of this sessions contain an expiration date. The `lifetime` of cookie depends on the activity of user. If he does not perform any actions within 10 minutes, the session ends. Otherwise, the session extends up to 10 minutes.
Also used UUID for session identification.

Instructions for user registration:

- Must ask for email
  - When the email is already taken return an error response.
- Must ask for username
- Must ask for password
  - The password must be encrypted when stored (this is a Bonus task)

The forum must checks if the email provided is present in the database and if all credentials are correct. It will check if the password is the same with the one provided and, if the password is not the same, it will return an error response.

#### Communication

In order for users to communicate between each other, they will have to be able to create posts and comments.

- Only registered users will be able to create posts and comments.
- When registered users are creating a post they can associate one or more categories to it.
  - The implementation and choice of the categories is up to you.
- The posts and comments should be visible to all users (registered or not).
- Non-registered users will only be able to see posts and comments.

#### Likes and Dislikes

Only registered users will be able to like or dislike posts and comments.

The number of likes and dislikes are visible by all users (registered or not).

#### Filter

Alslo, implemented a filter mechanism, that will allow users to filter the displayed posts by :

- categories
- created posts
- liked posts

You can look at filtering by categories as subforums. A subforum is a section of an online forum dedicated to a specific topic.

Note that the last two are only available for registered users and must refer to the logged in user.

#### Docker

For this project, also used Docker.



