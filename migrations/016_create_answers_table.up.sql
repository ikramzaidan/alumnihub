-- Create answers table
CREATE TABLE public.answers (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    form_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answer_text TEXT
);

-- Create indexes for answers
CREATE INDEX idx_answers_user_id ON public.answers(user_id);
CREATE INDEX idx_answers_form_id ON public.answers(form_id);
CREATE INDEX idx_answers_question_id ON public.answers(question_id);

-- Add foreign key constraints
ALTER TABLE public.answers
    ADD CONSTRAINT answers_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE public.answers
    ADD CONSTRAINT answers_form_id_fkey 
    FOREIGN KEY (form_id) REFERENCES public.forms(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE public.answers
    ADD CONSTRAINT answers_question_id_fkey 
    FOREIGN KEY (question_id) REFERENCES public.questions(id) ON UPDATE CASCADE ON DELETE CASCADE;