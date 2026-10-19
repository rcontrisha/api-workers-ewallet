-- PIN semua user: 123456  (bcrypt cost 10)
INSERT INTO account (id, fullname, phone, avatar, pin, is_verified) VALUES
  ('a0000000-0000-4000-8000-000000000001', 'Ucup', '628120000001', NULL, '$argon2id$v=19$m=65536,t=2,p=2$NQkWmNuwo6qfOz0mjIZI/w$801b4jZt9aHrt9Em8p9QKxrnFNkN41nbf0tfq7CXO7w', true),
  ('a0000000-0000-4000-8000-000000000002', 'Ucup 1',      '628120000002', NULL, '$argon2id$v=19$m=65536,t=2,p=2$NQkWmNuwo6qfOz0mjIZI/w$801b4jZt9aHrt9Em8p9QKxrnFNkN41nbf0tfq7CXO7w', true),
  ('a0000000-0000-4000-8000-000000000003', 'Ucup 2',      '628120000003', NULL, '$argon2id$v=19$m=65536,t=2,p=2$NQkWmNuwo6qfOz0mjIZI/w$801b4jZt9aHrt9Em8p9QKxrnFNkN41nbf0tfq7CXO7w', true),
  ('a0000000-0000-4000-8000-000000000004', 'Ucup 3',      '628120000004', NULL, '$argon2id$v=19$m=65536,t=2,p=2$NQkWmNuwo6qfOz0mjIZI/w$801b4jZt9aHrt9Em8p9QKxrnFNkN41nbf0tfq7CXO7w', true),
  ('a0000000-0000-4000-8000-000000000005', 'Ucup 4',      '628120000005', NULL, '$argon2id$v=19$m=65536,t=2,p=2$NQkWmNuwo6qfOz0mjIZI/w$801b4jZt9aHrt9Em8p9QKxrnFNkN41nbf0tfq7CXO7w', false),
  ('a0000000-0000-4000-8000-000000000006', 'Ucup 5',      '628120000006', NULL, '$argon2id$v=19$m=65536,t=2,p=2$NQkWmNuwo6qfOz0mjIZI/w$801b4jZt9aHrt9Em8p9QKxrnFNkN41nbf0tfq7CXO7w', false)
ON CONFLICT DO NOTHING;