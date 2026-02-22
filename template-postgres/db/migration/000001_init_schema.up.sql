CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar NOT NULL,
  email varchar NOT NULL UNIQUE,
  created_at timestamptz NOT NULL DEFAULT now()
);
