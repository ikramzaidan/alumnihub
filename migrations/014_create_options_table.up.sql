-- Create options table
CREATE TABLE public.options (
    id BIGSERIAL PRIMARY KEY,
    question_id INTEGER,
    option_text VARCHAR(255)
);

-- Create indexes for options
CREATE INDEX idx_options_question_id ON public.options(question_id);

-- Add foreign key constraint
ALTER TABLE public.options
    ADD CONSTRAINT options_question_id_fkey 
    FOREIGN KEY (question_id) REFERENCES public.questions(id) ON UPDATE CASCADE ON DELETE CASCADE;