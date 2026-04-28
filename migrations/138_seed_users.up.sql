-- Seed users data
INSERT INTO public.users (username, email, password, is_admin, created_at, updated_at) VALUES
('admin', 'admin@gmail.com', '$2a$12$qDysuB7aGhgtRCI08kP24OMVK3snloIpSRzhvbIBIusaGpdQ5vNIa', true, '2022-09-23 00:00:00', '2022-09-23 00:00:00')
ON CONFLICT (email) DO NOTHING;