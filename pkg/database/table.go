package database

func createTables() error {
	// SQL statement to create the users table
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        is_admin BOOLEAN DEFAULT FALSE
    );`

	// SQL statement to create the posts table
	createPostsTable := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	// SQL statement to create the comments table
	createCommentsTable := `
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER,
        user_id INTEGER,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	// SQL statement to create the likes table
	createLikesTable := `
    CREATE TABLE IF NOT EXISTS likes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        post_id INTEGER,
        comment_id INTEGER,
        is_like BOOLEAN,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (comment_id) REFERENCES comments(id),
        UNIQUE(user_id, post_id, comment_id)
    );`

	// SQL statement to create the categories table
	createCategoriesTable := `
    CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL
    );`

	// SQL statement to create the post_categories table (for many-to-many relationship)
	createPostCategoriesTable := `
    CREATE TABLE IF NOT EXISTS post_categories (
        post_id INTEGER,
        category_id INTEGER,
        PRIMARY KEY (post_id, category_id),
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (category_id) REFERENCES categories(id)
    );`

	// Slice containing all table creation SQL statements
	tables := []string{
		createUsersTable,
		createPostsTable,
		createCommentsTable,
		createLikesTable,
		createCategoriesTable,
		createPostCategoriesTable,
	}

	// Execute each table creation statement
	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return err
		}
	}

	return nil
}
