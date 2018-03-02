DROP TABLE budgets;
DROP TABLE users;

CREATE TABLE users(
  id serial PRIMARY KEY,
  username varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  encryptedPass text NOT NULL,
  lastAccess timestamp WITH TIME ZONE NOT NULL,
  verified boolean NOT NULL
);

CREATE TABLE budgets(
  id serial PRIMARY KEY,
  userId int references users(id),
  income int NOT NULL,
  rent int NOT NULL,
  wealth int NOT NULL
);
