CREATE TABLE projects (
    url TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT
);

INSERT INTO projects
VALUES
    ("https://localhost:5050/", "Portobloglio", "The website you're currently on!"),
    ("https://asuracomic.net/","Test", "Lorem ipsum dolor amet.");