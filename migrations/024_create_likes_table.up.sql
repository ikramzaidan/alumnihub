-- Create likes table
CREATE TABLE public.likes (
    id BIGSERIAL PRIMARY KEY,
    forum_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for likes
CREATE INDEX idx_likes_forum_id ON public.likes(forum_id);
CREATE INDEX idx_likes_user_id ON public.likes(user_id);
CREATE UNIQUE INDEX idx_likes_forum_user ON public.likes(forum_id, user_id);

-- Add foreign key constraints
ALTER TABLE public.likes
    ADD CONSTRAINT likes_forum_id_fkey 
    FOREIGN KEY (forum_id) REFERENCES public.forums(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE public.likes
    ADD CONSTRAINT likes_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;