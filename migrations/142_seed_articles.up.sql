-- Seed articles data
INSERT INTO public.articles (title, slug, body, status, created_at, updated_at, published_at) VALUES
('Kepala Sekolah Baru Resmi Dilantik', 'kepala-sekolah-baru-resmi-dilantik', '<p style="text-align:justify;"><strong>Lorem Ipsum</strong> is simply dummy text of the printing and typesetting industry...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-09-23 00:00:00', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),

('Reuni Akbar Alumni 2020', 'reuni-akbar-alumni-2020', '<p style="text-align:justify;">Lorem Ipsum is simply dummy text of the printing and typesetting industry...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-10-01 00:00:00', '2022-10-01 00:00:00', '2022-10-01 00:00:00'),

('Tips Karir untuk Alumni Baru', 'tips-karir-untuk-alumni-baru', '<p style="text-align:justify;">Lorem Ipsum is simply dummy text of the printing industry...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-10-05 00:00:00', '2022-10-05 00:00:00', '2022-10-05 00:00:00'),

('Alumni Berprestasi di Dunia Kerja', 'alumni-berprestasi-di-dunia-kerja', '<p style="text-align:justify;">Lorem Ipsum has been the industry standard dummy text...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-10-10 00:00:00', '2022-10-10 00:00:00', '2022-10-10 00:00:00'),

('Workshop Digital Marketing untuk Alumni', 'workshop-digital-marketing-untuk-alumni', '<p style="text-align:justify;">Lorem Ipsum is simply dummy text...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-10-15 00:00:00', '2022-10-15 00:00:00', '2022-10-15 00:00:00'),

('Program Beasiswa untuk Alumni', 'program-beasiswa-untuk-alumni', '<p style="text-align:justify;">Lorem Ipsum has survived not only five centuries...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-10-20 00:00:00', '2022-10-20 00:00:00', '2022-10-20 00:00:00'),

('Kegiatan Sosial Alumni 2022', 'kegiatan-sosial-alumni-2022', '<p style="text-align:justify;">Lorem Ipsum is simply dummy text of the printing...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-10-25 00:00:00', '2022-10-25 00:00:00', '2022-10-25 00:00:00'),

('Pelatihan Soft Skill Alumni', 'pelatihan-soft-skill-alumni', '<p style="text-align:justify;">Lorem Ipsum passages, and more recently...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-11-01 00:00:00', '2022-11-01 00:00:00', '2022-11-01 00:00:00'),

('Networking Event Alumni Nasional', 'networking-event-alumni-nasional', '<p style="text-align:justify;">Lorem Ipsum is simply dummy text...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-11-05 00:00:00', '2022-11-05 00:00:00', '2022-11-05 00:00:00'),

('Startup Alumni Sukses di Indonesia', 'startup-alumni-sukses-di-indonesia', '<p style="text-align:justify;">Lorem Ipsum has been the industry standard...</p><p style="text-align:justify;">&nbsp;</p>', 'published', '2022-11-10 00:00:00', '2022-11-10 00:00:00', '2022-11-10 00:00:00')

ON CONFLICT (slug) DO NOTHING;