-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR NOT NULL,
    content VARCHAR NOT NULL,
    user_id INTEGER,
    like_id INTEGER,
    dislike_id INTEGER,
    Comm_ID INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (like_id) REFERENCES PostLike(id),
    FOREIGN KEY (dislike_id) REFERENCES PostDislike(id),
    FOREIGN KEY (Comm_ID) REFERENCES Comments(comment_ID)
);

-- Create PostLike table
CREATE TABLE IF NOT EXISTS PostLike (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    likes INTEGER,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create PostDislike table
CREATE TABLE IF NOT EXISTS PostDislike (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dislikes INTEGER,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create CommentLike table
CREATE TABLE IF NOT EXISTS CommentLike (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    likes INTEGER,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create CommentDislike table
CREATE TABLE IF NOT EXISTS CommentDislike (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dislikes INTEGER,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create Comments table
CREATE TABLE IF NOT EXISTS Comments (
    comment_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    content VARCHAR NOT NULL,
    user_id INTEGER,
    like_id INTEGER,
    dislike_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (like_id) REFERENCES CommentLike(id),
    FOREIGN KEY (dislike_id) REFERENCES CommentDislike(id)
);
