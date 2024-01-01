CREATE TABLE IF NOT EXISTS auth_user (
    id VARCHAR(15) PRIMARY KEY,
    username VARCHAR(31) UNIQUE NOT NULL,
    email VARCHAR(50),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS user_key (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(15) NOT NULL,
    hashed_password VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES auth_user(id)
);
CREATE TABLE IF NOT EXISTS  user_session (
    id VARCHAR(127) PRIMARY KEY,
    user_id VARCHAR(15) NOT NULL,
    active_expires BIGINT NOT NULL,
    idle_expires BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES auth_user(id)
); 

CREATE TABLE IF NOT EXISTS template(
    id SERIAL PRIMARY KEY NOT NULL,
    template_name VARCHAR(40) NOT NULL,
    game_template JSONB NOT NULL,
    user_id VARCHAR(15) NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES auth_user(id)
);


CREATE TABLE IF NOT EXISTS saved_game (
    id SERIAL PRIMARY KEY NOT NULL,
    current_state JSONB NOT NULL,
    game_template JSONB NOT NULL,
    player1 VARCHAR(31) NOT NULL,
    player2 VARCHAR(31) NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (player1) REFERENCES auth_user(username),
    FOREIGN KEY (player2) REFERENCES auth_user(username)
);


