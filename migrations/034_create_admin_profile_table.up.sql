-- Create admin_profile table
CREATE TABLE public.admin_profile (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    is_super_admin BOOLEAN DEFAULT false,
    name VARCHAR(255),
    bio TEXT,
    sm_facebook VARCHAR(64),
    sm_instagram VARCHAR(64),
    sm_twitter VARCHAR(64),
    sm_tiktok VARCHAR(64)
);

-- Create indexes for admin_profile
CREATE INDEX idx_admin_profile_user_id ON public.admin_profile(user_id);
CREATE INDEX idx_admin_profile_is_super_admin ON public.admin_profile(is_super_admin);

-- Add foreign key constraint
ALTER TABLE public.admin_profile
    ADD CONSTRAINT admin_profile_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;