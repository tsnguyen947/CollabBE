DROP TABLE budgets;
DROP TABLE users;

CREATE TABLE users(
  id serial PRIMARY KEY,
  username varchar(255),
  email varchar(255),
  encryptedPass text,
  lastAccess timestamp WITH TIME ZONE,
  verified boolean
);

CREATE TABLE budgets(
  id serial PRIMARY KEY,
  userId int references users(id),
  income int,
  rent int,
  wealth int
);
