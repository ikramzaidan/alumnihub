-- Create alumni_jobs table
CREATE TABLE public.alumni_jobs (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    position VARCHAR(255),
    company VARCHAR(255),
    company_location VARCHAR(255),
    employment_type VARCHAR(255),
    start_year INTEGER,
    end_year INTEGER,
    currently_working BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for alumni_jobs
CREATE INDEX idx_alumni_jobs_user_id ON public.alumni_jobs(user_id);
CREATE INDEX idx_alumni_jobs_company ON public.alumni_jobs(company);
CREATE INDEX idx_alumni_jobs_currently_working ON public.alumni_jobs(currently_working);

-- Add foreign key constraint
ALTER TABLE public.alumni_jobs
    ADD CONSTRAINT alumni_jobs_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;