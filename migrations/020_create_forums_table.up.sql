-- Create forums table
CREATE TABLE public.forums (
    id BIGSERIAL PRIMARY KEY,
    forum_text TEXT,
    user_id INTEGER NOT NULL,
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for forums
CREATE INDEX idx_forums_user_id ON public.forums(user_id);
CREATE INDEX idx_forums_published_at ON public.forums(published_at);

-- Add foreign key constraint
ALTER TABLE public.forums
    ADD CONSTRAINT forums_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;