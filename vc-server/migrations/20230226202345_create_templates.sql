-- +goose Up
CREATE TABLE Templates (
    template_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(id),
    name VARCHAR(255) NOT NULL,
    starting_fen VARCHAR(255) NOT NULL,
);

CREATE TABLE Pieces (
    piece_id SERIAL PRIMARY KEY,
    name VARCHAR(2) NOT NULL,
    template_id INTEGER REFERENCES Templates(template_id)
);

CREATE TYPE Move_pattern AS (
  x INT,
  y INT
);

CREATE TABLE Move_patterns (
  id SERIAL PRIMARY KEY,
  type VARCHAR(10) CHECK (type IN ('slide', 'jump')),
  piece_id INTEGER REFERENCES Pieces(piece_id) ON DELETE CASCADE,
  template_id INTEGER REFERENCES Templates(template_id) ON DELETE CASCADE,
  move_pattern Move_pattern[]
);
