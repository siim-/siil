CREATE TABLE site_right (
	site_id CHAR(64) NOT NULL,
	user_code CHAR(11) NOT NULL,
	PRIMARY KEY (site_id, user_code),
	CONSTRAINT FK_site_right_site FOREIGN KEY (site_id) REFERENCES site (client_id) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT FK_site_right_user FOREIGN KEY (user_code) REFERENCES user (code) ON UPDATE CASCADE ON DELETE CASCADE
)