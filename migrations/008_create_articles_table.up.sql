-- Create articles table
CREATE TABLE public.articles (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(512),
    slug VARCHAR(255) UNIQUE NOT NULL,
    body TEXT,
    status public.article_status DEFAULT 'draft',
    image VARCHAR(255) DEFAULT 'public/no-image.png',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP
);

-- Create indexes for articles
CREATE INDEX idx_articles_slug ON public.articles(slug);
CREATE INDEX idx_articles_status ON public.articles(status);
CREATE INDEX idx_articles_published_at ON public.articles(published_at);
CREATE INDEX idx_articles_created_at ON public.articles(created_at);