--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(255) UNIQUE,
    email character varying(255) UNIQUE,
    password character varying(255),
    is_admin boolean,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


--
-- Name: alumni; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.alumni (
    id integer NOT NULL,
    nisn character varying(16) UNIQUE NOT NULL,
    nis character varying(16) UNIQUE NOT NULL,
    name character varying(512),
    gender character varying(1),
    phone character varying(16),
    graduation_year integer,
    class character varying(32),
    user_id integer UNIQUE DEFAULT NULL
);


--
-- Name: alumni; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.alumni_profile (
    id integer NOT NULL,
    user_id integer UNIQUE NOT NULL,
    alumni_id integer UNIQUE NOT NULL,
    bio text DEFAULT NULL,
    location character varying(255) DEFAULT NULL,
    sm_facebook character varying(64) DEFAULT NULL,
    sm_instagram character varying(64) DEFAULT NULL,
    sm_twitter character varying(64) DEFAULT NULL,
    sm_tiktok character varying(64) DEFAULT NULL
);


--
-- Name: articles; Type: TABLE; Schema: public; Owner: -
--

CREATE TYPE public.article_status AS ENUM ('draft', 'published');
CREATE TABLE public.articles (
    id integer NOT NULL,
    title character varying(512),
    slug character varying(255),
    body text,
    status public.article_status,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    published_at timestamp without time zone default NULL
);


--
-- Name: forms; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.forms (
    id integer NOT NULL,
    title character varying(255),
    description text,
    start_date timestamp without time zone,
    end_date timestamp without time zone,
    has_time_limit boolean,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


--
-- Name: questions; Type: TABLE; Schema: public; Owner: -
--

CREATE TYPE public.question_type AS ENUM ('multiple_choice', 'short_answer', 'long_answer');
CREATE TABLE public.questions (
    id integer NOT NULL,
    form_id integer,
    question_text text,
    type public.question_type,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


--
-- Name: options; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.options (
    id integer NOT NULL,
    question_id integer,
    option_text text
);


--
-- Name: answers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.answers (
    id integer NOT NULL,
    user_id integer NOT NULL,
    form_id integer NOT NULL,
    question_id integer NOT NULL,
    answer_text text
);


--
-- Name: forums; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.forums (
    id integer NOT NULL,
    forum_text text,
    user_id integer NOT NULL,
    published_at timestamp without time zone
);


--
-- Name: replies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.replies (
    id integer NOT NULL,
    forum_id integer NOT NULL,
    reply_text text,
    user_id integer NOT NULL,
    published_at timestamp without time zone
);


--
-- Name: likes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.likes (
    id integer NOT NULL,
    forum_id integer NOT NULL,
    user_id integer NOT NULL
);



--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: alumni_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.alumni ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.alumni_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: alumni_profile_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.alumni_profile ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.alumni_profile_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: articles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.articles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.articles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: forms_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.forms ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.forms_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: questions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.questions ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.questions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: options_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.options ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.options_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: answers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.answers ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.answers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: forums_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.forums ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.forums_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: replies_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.replies ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.replies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: likes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.likes ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.likes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: alumni alumni_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.alumni
    ADD CONSTRAINT alumni_pkey PRIMARY KEY (id);


--
-- Name: alumni alumni_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.alumni_profile
    ADD CONSTRAINT alumni_profile_pkey PRIMARY KEY (id);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id);


--
-- Name: forms forms_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.forms
    ADD CONSTRAINT forms_pkey PRIMARY KEY (id);


--
-- Name: questions questions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_pkey PRIMARY KEY (id);


--
-- Name: options options_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.options
    ADD CONSTRAINT options_pkey PRIMARY KEY (id);


--
-- Name: answers answers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.answers
    ADD CONSTRAINT answers_pkey PRIMARY KEY (id);


--
-- Name: forums forums_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_pkey PRIMARY KEY (id);


--
-- Name: replies replies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.replies
    ADD CONSTRAINT replies_pkey PRIMARY KEY (id);


--
-- Name: likes likes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_pkey PRIMARY KEY (id);


--
-- Name: alumni alumni_profile_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.alumni_profile
    ADD CONSTRAINT alumni_profile_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: alumni alumni_profile_alumni_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.alumni_profile
    ADD CONSTRAINT alumni_profile_alumni_id_fkey FOREIGN KEY (alumni_id) REFERENCES public.alumni(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: questions questions_form_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_form_id_fkey FOREIGN KEY (form_id) REFERENCES public.forms(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: options options_question_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.options
    ADD CONSTRAINT options_question_id_fkey FOREIGN KEY (question_id) REFERENCES public.questions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: answers answers_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.answers
    ADD CONSTRAINT answers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: answers answers_forms_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.answers
    ADD CONSTRAINT answers_form_id_fkey FOREIGN KEY (form_id) REFERENCES public.forms(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: answers answers_question_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.answers
    ADD CONSTRAINT answers_question_id_fkey FOREIGN KEY (question_id) REFERENCES public.questions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: forums forums_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: replies replies_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.replies
    ADD CONSTRAINT replies_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: replies replies_forum_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.replies
    ADD CONSTRAINT replies_forum_id_fkey FOREIGN KEY (forum_id) REFERENCES public.forums(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: likes likes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: likes likes_forum_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_forum_id_fkey FOREIGN KEY (forum_id) REFERENCES public.forums(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Data for Name: alumni; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.alumni (nisn, nis, name, gender, phone, graduation_year, class) FROM stdin;
0023978634	1202204216	Ikram Zaidan	M	081224939927	2020	XII-MIPA-1
0023978635	1202204217	Rayhan Ampurama	M	081123423323	2020	XII-MIPA-1
0023978636	1202204218	Alfatha Huga	M	098282828122	2020	XII-MIPA-1
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (username, email, password, is_admin, created_at, updated_at) FROM stdin;
admin	admin@gmail.com	$2a$12$qDysuB7aGhgtRCI08kP24OMVK3snloIpSRzhvbIBIusaGpdQ5vNIa	true	2022-09-23 00:00:00	2022-09-23 00:00:00
ikramzaidann	ikramzaidan820@gmail.com	$2a$12$qDysuB7aGhgtRCI08kP24OMVK3snloIpSRzhvbIBIusaGpdQ5vNIa	false	2022-09-23 00:00:00	2022-09-23 00:00:00
\.


--
-- Data for Name: articles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.articles (title, slug, body, status, created_at, updated_at, published_at) FROM stdin;
Kepala Sekolah Ganti	wawancara-aplikasi-alumni	<p style="text-align:justify;"><strong>Lorem Ipsum</strong> is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industrys standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</p><p style="text-align:justify;">&nbsp;</p><figure class="image"><img style="aspect-ratio:1200/1600;" src="http://localhost:8080/public/upload-3718730412.jpeg" width="1200" height="1600"></figure>	published	2022-09-23 00:00:00	2022-09-23 00:00:00	2022-09-23 00:00:00
\.


--
-- PostgreSQL database dump complete
--

