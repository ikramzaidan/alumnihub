-- Create alumni table
CREATE TABLE public.alumni (
    id BIGSERIAL PRIMARY KEY,
    nisn VARCHAR(16) UNIQUE NOT NULL,
    nis VARCHAR(16) UNIQUE NOT NULL,
    name VARCHAR(512),
    gender VARCHAR(1),
    phone VARCHAR(16),
    graduation_year INTEGER,
    class VARCHAR(32)
);

-- Create indexes for alumni
CREATE INDEX idx_alumni_nisn ON public.alumni(nisn);
CREATE INDEX idx_alumni_nis ON public.alumni(nis);
CREATE INDEX idx_alumni_graduation_year ON public.alumni(graduation_year);
CREATE INDEX idx_alumni_name ON public.alumni(name);