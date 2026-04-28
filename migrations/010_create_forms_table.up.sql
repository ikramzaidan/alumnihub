-- Create forms table
CREATE TABLE public.forms (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    has_time_limit BOOLEAN DEFAULT false,
    hidden BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for forms
CREATE INDEX idx_forms_title ON public.forms(title);
CREATE INDEX idx_forms_hidden ON public.forms(hidden);
CREATE INDEX idx_forms_start_date ON public.forms(start_date);
CREATE INDEX idx_forms_end_date ON public.forms(end_date);