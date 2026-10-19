INSERT INTO payment_method (id, gateway_payment_channel, name, logo, is_active) VALUES
  ('b0000000-0000-4000-8000-000000000001', 'VIRTUAL_ACCOUNT_BRI', 'Bank Rakyat Indonesia', '/logo/bri.png', true),
  ('b0000000-0000-4000-8000-000000000002', 'VIRTUAL_ACCOUNT_BCA', 'Bank Central Asia', '/logo/bca.png', true),
  ('b0000000-0000-4000-8000-000000000003', 'EMONEY_SHOPEE_PAY', 'ShopeePay', '/logo/spay.png', true),
  ('b0000000-0000-4000-8000-000000000004', 'EMONEY_DANA', 'Dana', '/logo/dana.png', true),
  ('b0000000-0000-4000-8000-000000000005', 'EMONEY_OVO', 'Ovo', '/logo/ovo.png', true)
ON CONFLICT DO NOTHING;