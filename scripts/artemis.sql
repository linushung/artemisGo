-- Create user
CREATE ROLE artemis WITH LOGIN CREATEDB CREATEROLE PASSWORD 'artemis';

-- Create database
CREATE DATABASE artemis WITH OWNER = artemis ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8' TABLESPACE = pg_default;

/* Ref:
1. https://www.postgresqltutorial.com/postgresql-data-types/
2. https://tapoueh.org/blog/2018/05/postgresql-data-types/
*/
CREATE TABLE poster (
    email VARCHAR(50) NOT NULL,
    username VARCHAR(20) NOT NULL,
    /* Ref: https://stackoverflow.com/questions/247304/what-data-type-to-use-for-hashed-password-field-and-what-length */
    password CHAR(60) NOT NULL,
    role VARCHAR(10) NOT NULL,
    image VARCHAR(100) DEFAULT '',
    bio VARCHAR(100) DEFAULT '',
    token VARCHAR(100) DEFAULT '',
    created_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    /* Ref:
    1. https://it.toolbox.com/question/the-difference-between-a-primary-key-and-a-surrogate-key-011407
    2. https://stackoverflow.com/questions/63090/surrogate-vs-natural-business-keys
    */
    PRIMARY KEY (email),
    UNIQUE (username),
    CHECK (role IN ('ADMIN', 'USER', 'VISITOR', 'UNKNOWN'))
);
CREATE INDEX username_index ON poster USING hash (username);

CREATE TABLE follower (
    email VARCHAR(50) NOT NULL,
    follower VARCHAR(20) NOT NULL,
    PRIMARY KEY (email, follower),
    FOREIGN KEY (email) REFERENCES poster (email) ON DELETE CASCADE
);

CREATE TABLE article (
    id UUID,
    slug VARCHAR(20) NOT NULL,
    title VARCHAR(20) NOT NULL,
    description VARCHAR(50) NOT NULL,
    body VARCHAR(200) NOT NULL,
    tagId SERIAL NOT NULL,
    favorite BOOLEAN DEFAULT False NOT NULL,
    favorite_count INTEGER DEFAULT 0,
    created_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE tag (
    id SERIAL,
    tag VARCHAR(15),
    PRIMARY KEY (id, tag),
    FOREIGN KEY (id) REFERENCES article (tagId) ON DELETE CASCADE
);

-- Set timezone for TIMESTAMPTZ column
SET TIMEZONE = 'Asia/Taipei';

-- Create function & triggers for auto update last_modified_time column in each table
--- This function sets any column named 'modified_time' to current timestamp for each row passed to it by the trigger
CREATE OR REPLACE FUNCTION update_modified_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_time = now();
    RETURN NEW;
END;
$$ language 'plpgsql';
--- Below triggers auto update 'modified_time' column in poster table to current timestamp
CREATE TRIGGER update_poster_modified BEFORE UPDATE ON poster FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
