-- Create questions table
CREATE TABLE public.questions (
    id BIGSERIAL PRIMARY KEY,
    form_id INTEGER,
    question_text TEXT,
    type public.question_type NOT NULL,
    extension BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for questions
CREATE INDEX idx_questions_form_id ON public.questions(form_id);
CREATE INDEX idx_questions_type ON public.questions(type);

-- Add foreign key constraint
ALTER TABLE public.questions
    ADD CONSTRAINT questions_form_id_fkey 
    FOREIGN KEY (form_id) REFERENCES public.forms(id) ON UPDATE CASCADE ON DELETE CASCADE;