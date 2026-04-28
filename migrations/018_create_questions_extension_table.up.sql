-- Create questions_extension table
CREATE TABLE public.questions_extension (
    id BIGSERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL,
    followup_question_id INTEGER NOT NULL,
    followup_option_value VARCHAR(255)
);

-- Create indexes for questions_extension
CREATE INDEX idx_questions_extension_question_id ON public.questions_extension(question_id);
CREATE INDEX idx_questions_extension_followup_question_id ON public.questions_extension(followup_question_id);

-- Add foreign key constraints
ALTER TABLE public.questions_extension
    ADD CONSTRAINT questions_extension_question_id_fkey 
    FOREIGN KEY (question_id) REFERENCES public.questions(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE public.questions_extension
    ADD CONSTRAINT questions_extension_followup_question_id_fkey 
    FOREIGN KEY (followup_question_id) REFERENCES public.questions(id) ON UPDATE CASCADE ON DELETE CASCADE;