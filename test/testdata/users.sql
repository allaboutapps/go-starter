--
-- PostgreSQL database dump
--
-- Dumped from database version 12.4
-- Dumped by pg_dump version 12.4
SET statement_timeout = 0;

SET lock_timeout = 0;

SET idle_in_transaction_session_timeout = 0;

SET client_encoding = 'UTF8';

SET standard_conforming_strings = ON;

SELECT
    pg_catalog.set_config('search_path', '', FALSE);

SET check_function_bodies = FALSE;

SET xmloption = content;

SET client_min_messages = warning;

SET row_security = OFF;

ALTER TABLE IF EXISTS ONLY public.users
    DROP CONSTRAINT IF EXISTS users_username_key;

ALTER TABLE IF EXISTS ONLY public.users
    DROP CONSTRAINT IF EXISTS users_pkey;

DROP TABLE IF EXISTS public.users;

DROP EXTENSION IF EXISTS "uuid-ossp";

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner:
--
COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: dbuser
--
CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4 () NOT NULL,
    username character varying(255),
    password text,
    is_active boolean NOT NULL,
    scopes text[] NOT NULL,
    last_authenticated_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

ALTER TABLE public.users OWNER TO dbuser;

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: dbuser
--
COPY public.users (id, username, password, is_active, scopes, last_authenticated_at, created_at, updated_at) FROM stdin;
44a4b372-9d45-42a7-a4bd-9ba78b580e09	test_user1@example.com	$argon2id$v=19$m=65536,t=1,p=4$bM/6DaUMUQlr8CPYHOcFwA$Cq0p8d3IKEy1+G3VfHRXDRdc15wHJtTx+UGsXD6bbWY	t	{app}	2020-09-16 14:39:34.098333+00	2020-09-16 12:37:27.034337+00	2020-09-16 14:39:34.098334+00
738a3a72-2267-4727-beef-44193491c7d0	test_user2@example.com	$argon2id$v=19$m=65536,t=1,p=4$goeC9sTUXfQMpLbdqHbIiA$nadmsRl8d0HTGuKOmjg1WGGOfvJkfPUb8aSj48t7upk	t	{app}	2020-11-10 13:28:11.259994+00	2020-11-10 13:28:11.259998+00	2020-11-10 13:28:11.259998+00
8c39db83-c355-4e55-b10a-ac9bf0b15e50	test_user3@example.com	$argon2id$v=19$m=65536,t=1,p=4$CWvCxZhldc/UmlZaHje7jg$igRKSJjHxlUR/5l21FX1aQIGbm+1J30/L1fLQGduy/U	t	{app}	2021-02-24 13:57:40.870073+00	2021-02-24 13:57:40.870076+00	2021-02-24 13:57:40.870076+00
\.

--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: dbuser
--
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);

--
-- PostgreSQL database dump complete
--
