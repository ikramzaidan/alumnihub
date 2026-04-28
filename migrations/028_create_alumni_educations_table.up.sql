-- Create alumni_educations table
CREATE TABLE public.alumni_educations (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    school_name VARCHAR(255),
    school_degree VARCHAR(255),
    school_study_major VARCHAR(255),
    start_year INTEGER,
    end_year INTEGER,
    currently_studying BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for alumni_educations
CREATE INDEX idx_alumni_educations_user_id ON public.alumni_educations(user_id);
CREATE INDEX idx_alumni_educations_school_name ON public.alumni_educations(school_name);
CREATE INDEX idx_alumni_educations_currently_studying ON public.alumni_educations(currently_studying);

-- Add foreign key constraint
ALTER TABLE public.alumni_educations
    ADD CONSTRAINT alumni_educations_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;