-- Seed admin_profile data
INSERT INTO public.admin_profile (user_id, is_super_admin, name, bio) VALUES
(1, true, 'Admin', 'Admin portal alumni SMA Telkom Bandung')
ON CONFLICT (user_id) DO NOTHING;