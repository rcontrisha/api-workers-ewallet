-- Topup: amount = saldo yang masuk, total = amount + tax (yang dibayar user).
INSERT INTO transaction
  (trx_code, type, status, amount, tax, discount, total,
   sender_wallet_id, receiver_wallet_id, payment_method_id, note, created_at)
VALUES
  ('TRX-SEED-001', 'topup', 'success', 500000, 50000, 0, 550000, NULL,
     'a0000000-0000-4000-8000-000000000001', 'b0000000-0000-4000-8000-000000000002', NULL, now() - interval '6 days 3 hours'),
  ('TRX-SEED-002', 'topup', 'success', 200000, 20000, 0, 220000, NULL,
     'a0000000-0000-4000-8000-000000000003', 'b0000000-0000-4000-8000-000000000003', NULL, now() - interval '6 days 1 hour'),
  ('TRX-SEED-003', 'topup', 'success', 100000, 10000, 0, 110000, NULL,
     'a0000000-0000-4000-8000-000000000004', 'b0000000-0000-4000-8000-000000000001', NULL, now() - interval '5 days 6 hours'),
  ('TRX-SEED-013', 'topup', 'failed',  100000, 10000, 0, 110000, NULL,
     'a0000000-0000-4000-8000-000000000001', 'b0000000-0000-4000-8000-000000000003', NULL, now() - interval '1 day 5 hours')
ON CONFLICT DO NOTHING;

INSERT INTO transaction
  (trx_code, type, status, amount, tax, discount, total,
   sender_wallet_id, receiver_wallet_id, note, deleted_by_sender, created_at)
VALUES
  ('TRX-SEED-004', 'transfer', 'success',  50000, 0, 0,  50000,
     'a0000000-0000-4000-8000-000000000001', 'a0000000-0000-4000-8000-000000000002', 'Bayar kopi',   false, now() - interval '5 days 2 hours'),
  ('TRX-SEED-005', 'transfer', 'success',  75000, 0, 0,  75000,
     'a0000000-0000-4000-8000-000000000003', 'a0000000-0000-4000-8000-000000000001', 'Patungan',     false, now() - interval '4 days 8 hours'),
  ('TRX-SEED-006', 'transfer', 'success',  25000, 0, 0,  25000,
     'a0000000-0000-4000-8000-000000000001', 'a0000000-0000-4000-8000-000000000004', NULL,           false, now() - interval '4 days 2 hours'),
  ('TRX-SEED-007', 'transfer', 'success',  40000, 0, 0,  40000,
     'a0000000-0000-4000-8000-000000000004', 'a0000000-0000-4000-8000-000000000001', 'Ganti bensin', false, now() - interval '3 days 7 hours'),
  ('TRX-SEED-008', 'transfer', 'success', 100000, 0, 0, 100000,
     'a0000000-0000-4000-8000-000000000001', 'a0000000-0000-4000-8000-000000000005', 'Titip beli',   false, now() - interval '3 days 1 hour'),
  ('TRX-SEED-009', 'transfer', 'success',  20000, 0, 0,  20000,
     'a0000000-0000-4000-8000-000000000002', 'a0000000-0000-4000-8000-000000000001', NULL,           false, now() - interval '2 days 6 hours'),
  ('TRX-SEED-010', 'transfer', 'success',  60000, 0, 0,  60000,
     'a0000000-0000-4000-8000-000000000001', 'a0000000-0000-4000-8000-000000000003', 'Makan siang',  false, now() - interval '2 days 1 hour'),
  ('TRX-SEED-011', 'transfer', 'success',  30000, 0, 0,  30000,
     'a0000000-0000-4000-8000-000000000003', 'a0000000-0000-4000-8000-000000000006', NULL,           false, now() - interval '1 day 8 hours'),
  -- contoh soft delete: sudah dihapus dari history pengirim, masih tampil di penerima
  ('TRX-SEED-012', 'transfer', 'success',  15000, 0, 0,  15000,
     'a0000000-0000-4000-8000-000000000001', 'a0000000-0000-4000-8000-000000000002', 'Parkir',       true,  now() - interval '1 day 2 hours')
ON CONFLICT DO NOTHING;