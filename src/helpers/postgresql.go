/*
author Vyacheslav Kryuchenko
*/
package helpers

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	SCHEMA = `
CREATE TABLE IF NOT EXISTS users
(
  id     SERIAL                                NOT NULL
    CONSTRAINT users_pkey
    PRIMARY KEY,
  login  VARCHAR(32)                           NOT NULL,
  email  VARCHAR(64),
  active BOOLEAN DEFAULT TRUE                  NOT NULL,
  admin  BOOLEAN DEFAULT FALSE                 NOT NULL,
  quota  BIGINT DEFAULT '5368709120' :: BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex
  ON users (login);

CREATE TABLE IF NOT EXISTS history
(
  action_time TIMESTAMP DEFAULT now() NOT NULL
    CONSTRAINT history_pkey
    PRIMARY KEY,
  user_id     INTEGER                 NOT NULL
    CONSTRAINT history_users__fk
    REFERENCES users
    ON UPDATE CASCADE ON DELETE CASCADE,
  command     TEXT                    NOT NULL,
  success     BOOLEAN DEFAULT TRUE    NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions
(
  user_id    INTEGER                                            NOT NULL
    CONSTRAINT sessions_user_id_pk
    PRIMARY KEY
    CONSTRAINT sessions_users__fk
    REFERENCES users
    ON UPDATE CASCADE ON DELETE CASCADE,
  expire     TIMESTAMP DEFAULT (now() + '08:00:00' :: INTERVAL) NOT NULL,
  session_id UUID                                               NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS sessions_user_id_uindex
  ON sessions (user_id);

CREATE TABLE IF NOT EXISTS container
(
  user_id INTEGER          NOT NULL
    CONSTRAINT data_size_users_id_fk
    REFERENCES users
    ON UPDATE CASCADE ON DELETE CASCADE,
  name    VARCHAR(64)      NOT NULL
    CONSTRAINT data_size_pkey
    PRIMARY KEY,
  size    BIGINT DEFAULT 0 NOT NULL,
  port    SMALLINT         NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS data_size_container_uindex
  ON container (name);

CREATE UNIQUE INDEX IF NOT EXISTS container_port_uindex
  ON container (port);
`
)

type PostgresConfig struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	DBName  string `json:"dbname"`
	SSL     string `json:"ssl,omitempty"`
}

func (pg *PostgresConfig) Connect() (*sql.DB, error) {
	var ssl string
	if pg.SSL != "" {
		ssl = pg.SSL
	} else {
		ssl = "disable"
	}
	connectionParams := fmt.Sprintf("postgres://%s@%s/%s?sslmode=%s", pg.Auth, pg.Address, pg.DBName, ssl)
	return sql.Open("postgres", connectionParams)
}

func (pg *PostgresConfig) InitDatabase() error {
	sqlCon, err := pg.Connect()
	if err != nil {
		return err
	}
	defer sqlCon.Close()
	_, err = sqlCon.Exec(SCHEMA)
	if err != nil {
		return err
	}
	return nil
}
