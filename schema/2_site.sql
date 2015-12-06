DROP TABLE IF EXISTS site;
CREATE TABLE site (
	client_id CHAR(64) NOT NULL,
	private_id CHAR(128) NOT NULL,
	owner INT UNSIGNED NOT NULL,
	name VARCHAR(40) NOT NULL,
	domain VARCHAR(255) NOT NULL,
	callback_url TEXT NOT NULL,
	cancel_url TEXT NOT NULL,
	PRIMARY KEY (client_id),
	CONSTRAINT FK_site_owner FOREIGN KEY (owner) REFERENCES user (id) ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 DEFAULT COLLATE utf8_unicode_ci;
INSERT INTO site (client_id, private_id, owner, name, domain, callback_url, cancel_url) VALUES ('a1s2d34', 'trakyll', 1, 'Siil.lan', 'siil.lan', 'https://siil.lan/success', 'https://siil.lan/fail');