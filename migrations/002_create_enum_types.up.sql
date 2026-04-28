-- Create enum types
CREATE TYPE public.article_status AS ENUM ('draft', 'published');
CREATE TYPE public.question_type AS ENUM ('multiple_choice', 'short_answer', 'long_answer');