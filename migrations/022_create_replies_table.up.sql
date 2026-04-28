-- Create replies table
CREATE TABLE public.replies (
    id BIGSERIAL PRIMARY KEY,
    forum_id INTEGER NOT NULL,
    reply_text TEXT,
    user_id INTEGER NOT NULL,
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for replies
CREATE INDEX idx_replies_forum_id ON public.replies(forum_id);
CREATE INDEX idx_replies_user_id ON public.replies(user_id);

-- Add foreign key constraints
ALTER TABLE public.replies
    ADD CONSTRAINT replies_forum_id_fkey 
    FOREIGN KEY (forum_id) REFERENCES public.forums(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE public.replies
    ADD CONSTRAINT replies_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;