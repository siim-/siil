DROP TABLE IF EXISTS user;
CREATE TABLE user (
	id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
	code CHAR(11) NOT NULL,
	first_name VARCHAR(60) NOT NULL,
	last_name VARCHAR(60) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE UNQ_user_code (code)
);

INSERT INTO user (code, first_name, last_name) VALUES ('12345678987', 'Tester', 'Ester');