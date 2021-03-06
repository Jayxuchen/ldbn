# Create table statements for an SQLite3 database

DROP TABLE IF EXISTS assignment;
CREATE TABLE assignment (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL, 
  name VARCHAR NOT NULL,
  modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  xml TEXT NOT NULL
);

DROP TABLE IF EXISTS user;
CREATE TABLE user (
  user_id INTEGER PRIMARY KEY,
  name VARCHAR NOT NULL,
  pass_hash VARCHAR NOT NULL,
  pass_salt VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  is_active INTEGER NOT NULL DEFAULT 0,
  is_admin INTEGER NOT NULL DEFAULT 0,
  is_su INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS session;
CREATE TABLE session (
  session_string VARCHAR PRIMARY KEY,
  user_id INTEGER NOT NULL,
  added_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) WITHOUT ROWID;


DROP TABLE IF EXISTS activation;
CREATE TABLE activation (
  activation_string VARCHAR PRIMARY KEY,
  user_id INTEGER NOT NULL,
  added_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) WITHOUT ROWID;

DROP TABLE IF EXISTS comment;
CREATE TABLE comment (
  id INTEGER PRIMARY KEY,
  assignment_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  comment_val TEXT NOT NULL,
  modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

