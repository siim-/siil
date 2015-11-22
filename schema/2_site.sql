CREATE TABLE site (
	client_id CHAR(64) NOT NULL,
	private_id CHAR(128) NOT NULL,
	owner CHAR(11) NOT NULL,
	name VARCHAR(40) NOT NULL,
	domain VARCHAR(255) NOT NULL,
	callback_url TEXT NOT NULL,
	cancel_url TEXT NOT NULL,
	PRIMARY KEY (client_id),
	CONSTRAINT FK_site_owner FOREIGN KEY (owner) REFERENCES user (code) ON UPDATE CASCADE ON DELETE CASCADE
);