-- Create alumni_profile table
CREATE TABLE public.alumni_profile (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    alumni_id INTEGER UNIQUE NOT NULL,
    bio TEXT,
    location VARCHAR(255),
    sm_facebook VARCHAR(64),
    sm_instagram VARCHAR(64),
    sm_twitter VARCHAR(64),
    sm_tiktok VARCHAR(64)
);

-- Create indexes for alumni_profile
CREATE INDEX idx_alumni_profile_user_id ON public.alumni_profile(user_id);
CREATE INDEX idx_alumni_profile_alumni_id ON public.alumni_profile(alumni_id);

-- Add foreign key constraints
ALTER TABLE public.alumni_profile
    ADD CONSTRAINT alumni_profile_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE public.alumni_profile
    ADD CONSTRAINT alumni_profile_alumni_id_fkey 
    FOREIGN KEY (alumni_id) REFERENCES public.alumni(id) ON UPDATE CASCADE ON DELETE CASCADE;