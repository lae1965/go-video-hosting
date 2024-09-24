CREATE TYPE ROLE AS ENUM ('userRole', 'adminRole', 'superAdminRole');

CREATE TABLE USERS (
  id serial PRIMARY KEY,
  nickName VARCHAR(64) NOT NULL UNIQUE,
  email VARCHAR(64) NOT NULL UNIQUE,
  passwordHash VARCHAR(64) NOT NULL,
  firstName VARCHAR(64) NULL,
  lastName VARCHAR(64) NULL,
  birthDate DATE NULL,
  role ROLE DEFAULT 'userRole',
  activateLink VARCHAR(64) NULL,
  isActivate BOOLEAN DEFAULT false,
  isBanned BOOLEAN DEFAULT false,
  channelsCount INTEGER DEFAULT 0,
  createTimestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

