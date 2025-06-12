CREATE TABLE projects (
    url TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT
);

INSERT INTO projects
VALUES
    ("http://localhost:5050/", "Portobloglio", "The website you're currently on!"),
    ("https://asuracomic.net/","Test", "Lorem ipsum dolor amet.");

CREATE TABLE blogs (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO categories (id, name) VALUES
  ('cat-dev', 'Development'),
  ('cat-life', 'Life & Thoughts'),
  ('cat-howto', 'How-To Guides');

INSERT INTO tags (id, name) VALUES
  ('tag-go', 'golang'),
  ('tag-web', 'webdev'),
  ('tag-htmx', 'htmx'),
  ('tag-sql', 'sqlite'),
  ('tag-life', 'philosophy');

INSERT INTO blogs (id, title, content, category_id, created_at, updated_at) VALUES
(
  'blog-001',
  'Getting Started with Go',
  '# Getting Started with Go\n\nGo is a statically typed, compiled language designed for simplicity and performance.\n\n## Why Go?\n- Simple syntax\n- Great performance\n- Built-in concurrency\n\nLearn more at [golang.org](https://golang.org)',
  'cat-dev',
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
),
(
  'blog-002',
  'Building Reactive UIs with HTMX',
  '# HTMX and the Modern Web\n\nHTMX allows you to build reactive web interfaces with minimal JavaScript.\n\n```html\n<button hx-get="/hello" hx-target="#output">Click me</button>\n```\n\nIt’s lightweight and powerful.',
  'cat-howto',
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
),
(
  'blog-003',
  'Things I Learned the Hard Way',
  '# Lessons from Experience\n\n> “Success is not final, failure is not fatal.”\n\nHere are 3 things I wish I knew earlier:\n\n1. Document everything.\n2. Keep things simple.\n3. You don’t need every feature.\n\nStay humble. Stay curious.',
  'cat-life',
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
);
