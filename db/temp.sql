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
    id BLOB PRIMARY KEY,
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
  X'0197C0955C76733193E59B05044C8B50',
  'Getting Started with Go',
  '# Getting Started with Go\n\nGo is a statically typed, compiled language designed at Google.\n\n```go\nfmt.Println("Hello, World!")\n```\n\nRead more at [golang.org](https://golang.org)',
  'cat-dev',
  '2025-06-01 10:00:00',
  '2025-06-01 10:00:00'
),
(
  X'0197C0955C7670A8A6B208A2D7F28D59',
  'Reflections on a Quiet Morning',
  '# Reflections on a Quiet Morning\n\nSometimes, all we need is a cup of coffee and silence.\n\n> "In stillness, the world resets."\n\nHere’s what I learned this week...',
  'cat-life',
  '2025-06-02 08:45:00',
  '2025-06-02 08:45:00'
),
(
  X'0197C0955C767800B9E86AC3235D4AFA',
  'Deploying SQLite in Production',
  '# Deploying SQLite in Production\n\nIs it ever a good idea?\n\n- ✅ Lightweight\n- ✅ Embedded\n- ⚠️ Limited concurrency\n\nLearn how to mitigate issues with proper caching.',
  'cat-dev',
  '2025-06-03 13:20:00',
  '2025-06-03 13:20:00'
),
(
  X'0197C0955C767B12B0A7594B0C69F02B',
  'How to Bake Bread at Home',
  '# How to Bake Bread at Home\n\n**Ingredients**:\n- 3 cups of flour\n- 1 tsp yeast\n- 1 cup water\n\n**Steps**:\n1. Mix\n2. Knead\n3. Bake\n\nEnjoy your fresh bread!',
  'cat-howto',
  '2025-06-04 16:05:00',
  '2025-06-04 16:05:00'
),
(
  X'0197C0955C767F85A8F402A10FAD78E1',
  'Why I Switched to Markdown Journaling',
  '# Why I Switched to Markdown Journaling\n\nUsing plain text has helped me:\n\n- Stay consistent\n- Avoid distractions\n- Sync across devices easily\n\nHere’s my daily template...',
  'cat-life',
  '2025-06-05 07:50:00',
  '2025-06-05 07:50:00'
);
