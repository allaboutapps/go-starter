--
-- PostgreSQL database dump
--
-- Dumped from database version 12.4
-- Dumped by pg_dump version 12.10 (Debian 12.10-1.pgdg100+1)
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

ALTER TABLE IF EXISTS ONLY public.gorp_migrations
    DROP CONSTRAINT IF EXISTS gorp_migrations_pkey;

DROP TABLE IF EXISTS public.gorp_migrations;

DROP EXTENSION IF EXISTS "uuid-ossp";

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner:
--
COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';

--
-- Name: gorp_migrations; Type: TABLE; Schema: public; Owner: dbuser
--
CREATE TABLE public.gorp_migrations (
    id text NOT NULL,
    applied_at timestamp with time zone
);

ALTER TABLE public.gorp_migrations OWNER TO dbuser;

--
-- Data for Name: gorp_migrations; Type: TABLE DATA; Schema: public; Owner: dbuser
--
COPY public.gorp_migrations (id, applied_at) FROM stdin;
20200428064736-install-extension-uuid.sql	2022-03-28 17:04:00.580862+00
\.

--
-- Name: gorp_migrations gorp_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--
ALTER TABLE ONLY public.gorp_migrations
    ADD CONSTRAINT gorp_migrations_pkey PRIMARY KEY (id);

--
-- PostgreSQL database dump complete
--
