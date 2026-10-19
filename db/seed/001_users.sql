-- Password semua user: Password123!  (bcrypt cost 10)
INSERT INTO users (id, email, password) VALUES
  ('a0000000-0000-4000-8000-000000000001', 'ucup@example.com',  '$argon2id$v=19$m=65536,t=2,p=2$lwIJnVVF2rxx3cliC3Kktw$sKcBUwKmZATi2MpSiEmvQaCeDIGOj0MXyAF2BVS0bx4'),
  ('a0000000-0000-4000-8000-000000000002', 'ucup1@example.com', '$argon2id$v=19$m=65536,t=2,p=2$lwIJnVVF2rxx3cliC3Kktw$sKcBUwKmZATi2MpSiEmvQaCeDIGOj0MXyAF2BVS0bx4'),
  ('a0000000-0000-4000-8000-000000000003', 'ucup2@example.com', '$argon2id$v=19$m=65536,t=2,p=2$lwIJnVVF2rxx3cliC3Kktw$sKcBUwKmZATi2MpSiEmvQaCeDIGOj0MXyAF2BVS0bx4'),
  ('a0000000-0000-4000-8000-000000000004', 'ucup3@example.com', '$argon2id$v=19$m=65536,t=2,p=2$lwIJnVVF2rxx3cliC3Kktw$sKcBUwKmZATi2MpSiEmvQaCeDIGOj0MXyAF2BVS0bx4'),
  ('a0000000-0000-4000-8000-000000000005', 'ucup4@example.com', '$argon2id$v=19$m=65536,t=2,p=2$lwIJnVVF2rxx3cliC3Kktw$sKcBUwKmZATi2MpSiEmvQaCeDIGOj0MXyAF2BVS0bx4'),
  ('a0000000-0000-4000-8000-000000000006', 'ucup5@example.com', '$argon2id$v=19$m=65536,t=2,p=2$lwIJnVVF2rxx3cliC3Kktw$sKcBUwKmZATi2MpSiEmvQaCeDIGOj0MXyAF2BVS0bx4')
ON CONFLICT DO NOTHING;