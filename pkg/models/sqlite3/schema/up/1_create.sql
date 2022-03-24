CREATE TABLE IF NOT EXISTS user (
			user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(20) NOT NULL UNIQUE,
			email VARCHAR(30) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			registration_date DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS user_session (
			uuid TEXT,
			expires DATETIME,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id)
				REFERENCES user (user_id)
				ON UPDATE CASCADE
				ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS post (
			post_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(50) NOT NULL CHECK(title != ''),
			content TEXT NOT NULL CHECK(content != ''),
			creation_date DATETIME NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id)
				REFERENCES user (user_id)
				ON UPDATE CASCADE
				ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS comment (
			comment_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			msg  TEXT NOT NULL CHECK(msg != ''),
			creation_date DATETIME NOT NULL,
			user_id INTEGER NOT NULL REFERENCES user (user_id),
			post_id INTEGER NOT NULL REFERENCES post (post_id)
				ON UPDATE CASCADE
				ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS category (
			category_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			category_name VARCHAR(50) UNIQUE
		);
		
		CREATE TABLE IF NOT EXISTS post_category (
			post_id INTEGER NOT NULL REFERENCES post (post_id) ON DELETE CASCADE,
			category_id INTEGER NOT NULL REFERENCES category (category_id)
				ON UPDATE CASCADE
				ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS rating_info (
			user_id INTEGER NOT NULL REFERENCES user (user_id),
			post_id INTEGER NOT NULL REFERENCES post (post_id),
			reaction VARCHAR(10) NOT NULL,
			CONSTRAINT UC_rating_info UNIQUE (user_id, post_id)	
		);

		CREATE TABLE IF NOT EXISTS comments_rating (
			user_id INTEGER NOT NULL REFERENCES user (user_id),
			comment_id INTEGER NOT NULL REFERENCES comment (comment_id),
			reaction VARCHAR(10) NOT NULL,
			CONSTRAINT UC_comment_rating UNIQUE (user_id, comment_id)
		)