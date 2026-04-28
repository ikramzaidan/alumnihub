-- Create jobs table
CREATE TABLE public.jobs (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    job_position VARCHAR(255),
    company VARCHAR(255),
    job_location VARCHAR(255),
    job_type VARCHAR(255),
    min_salary INTEGER,
    max_salary INTEGER,
    description TEXT,
    closed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for jobs
CREATE INDEX idx_jobs_user_id ON public.jobs(user_id);
CREATE INDEX idx_jobs_company ON public.jobs(company);
CREATE INDEX idx_jobs_job_location ON public.jobs(job_location);
CREATE INDEX idx_jobs_closed ON public.jobs(closed);
CREATE INDEX idx_jobs_created_at ON public.jobs(created_at);

-- Add foreign key constraint
ALTER TABLE public.jobs
    ADD CONSTRAINT jobs_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;