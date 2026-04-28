-- Seed alumni data
INSERT INTO public.alumni (nisn, nis, name, gender, phone, graduation_year, class) VALUES
('0023978634', '1202204216', 'Ikram Zaidan', 'M', '081224939927', 2020, 'XII-MIPA-1'),
('0023978635', '1202204217', 'Rayhan Ampurama', 'M', '081123423323', 2020, 'XII-MIPA-1'),
('0023978636', '1202204218', 'Alfatha Huga', 'M', '098282828122', 2020, 'XII-MIPA-1')
ON CONFLICT (nisn) DO NOTHING;