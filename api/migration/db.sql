--
-- PostgreSQL database dump
--

-- Dumped from database version 16.0 (Homebrew)
-- Dumped by pg_dump version 16.0 (Homebrew)

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
-- Name: customers; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.customers (
    id integer NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.customers OWNER TO labyrinth;

--
-- Name: customers_id_seq; Type: SEQUENCE; Schema: public; Owner: labyrinth
--

ALTER TABLE public.customers ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.customers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: orders; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    widget_id integer NOT NULL,
    transaction_id integer NOT NULL,
    status_id integer NOT NULL,
    quantity integer NOT NULL,
    amount integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    customers_id integer NOT NULL
);


ALTER TABLE public.orders OWNER TO labyrinth;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: labyrinth
--

ALTER TABLE public.orders ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO labyrinth;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.sessions (
    token text NOT NULL,
    data bytea NOT NULL,
    expiry timestamp with time zone NOT NULL
);


ALTER TABLE public.sessions OWNER TO labyrinth;

--
-- Name: statuses; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.statuses (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.statuses OWNER TO labyrinth;

--
-- Name: statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: labyrinth
--

ALTER TABLE public.statuses ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.statuses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.tokens (
    id integer NOT NULL,
    user_id integer,
    name character varying(255),
    email character varying(255) NOT NULL,
    updated_at timestamp without time zone DEFAULT now(),
    created_at timestamp without time zone DEFAULT now(),
    token_hash bytea,
    expiry timestamp without time zone
);


ALTER TABLE public.tokens OWNER TO labyrinth;

--
-- Name: tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: labyrinth
--

CREATE SEQUENCE public.tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tokens_id_seq OWNER TO labyrinth;

--
-- Name: tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: labyrinth
--

ALTER SEQUENCE public.tokens_id_seq OWNED BY public.tokens.id;


--
-- Name: transaction_statuses; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.transaction_statuses (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.transaction_statuses OWNER TO labyrinth;

--
-- Name: transaction_statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: labyrinth
--

ALTER TABLE public.transaction_statuses ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.transaction_statuses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.transactions (
    id integer NOT NULL,
    amount integer NOT NULL,
    currency character varying(255) NOT NULL,
    last_four character varying(255) NOT NULL,
    bank_return_code character varying(255) NOT NULL,
    transaction_status_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expiry_month integer DEFAULT 0 NOT NULL,
    expiry_year integer DEFAULT 0 NOT NULL,
    payment_intent character varying(255) DEFAULT ''::character varying NOT NULL,
    payment_method character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.transactions OWNER TO labyrinth;

--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: labyrinth
--

ALTER TABLE public.transactions ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.users OWNER TO labyrinth;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: labyrinth
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
-- Name: widgets; Type: TABLE; Schema: public; Owner: labyrinth
--

CREATE TABLE public.widgets (
    id integer NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL,
    description text DEFAULT ''::text NOT NULL,
    inventory_level integer NOT NULL,
    price integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    image character varying(255) DEFAULT ''::character varying NOT NULL,
    is_recurring boolean DEFAULT false NOT NULL,
    plan_id character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.widgets OWNER TO labyrinth;

--
-- Name: widgets_id_seq; Type: SEQUENCE; Schema: public; Owner: labyrinth
--

ALTER TABLE public.widgets ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.widgets_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: tokens id; Type: DEFAULT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.tokens ALTER COLUMN id SET DEFAULT nextval('public.tokens_id_seq'::regclass);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (token);


--
-- Name: statuses statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.statuses
    ADD CONSTRAINT statuses_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: transaction_statuses transaction_statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.transaction_statuses
    ADD CONSTRAINT transaction_statuses_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: widgets widgets_pkey; Type: CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.widgets
    ADD CONSTRAINT widgets_pkey PRIMARY KEY (id);


--
-- Name: orders_customers_id_fk; Type: INDEX; Schema: public; Owner: labyrinth
--

CREATE INDEX orders_customers_id_fk ON public.orders USING btree (customers_id);


--
-- Name: orders_statuses_id_fk; Type: INDEX; Schema: public; Owner: labyrinth
--

CREATE INDEX orders_statuses_id_fk ON public.orders USING btree (status_id);


--
-- Name: orders_transactions_id_fk; Type: INDEX; Schema: public; Owner: labyrinth
--

CREATE INDEX orders_transactions_id_fk ON public.orders USING btree (transaction_id);


--
-- Name: orders_widgets_id_fk; Type: INDEX; Schema: public; Owner: labyrinth
--

CREATE INDEX orders_widgets_id_fk ON public.orders USING btree (widget_id);


--
-- Name: sessions_expiry_idx; Type: INDEX; Schema: public; Owner: labyrinth
--

CREATE INDEX sessions_expiry_idx ON public.sessions USING btree (expiry);


--
-- Name: transactions_transaction_statuses_id_fk; Type: INDEX; Schema: public; Owner: labyrinth
--

CREATE INDEX transactions_transaction_statuses_id_fk ON public.transactions USING btree (transaction_status_id);


--
-- Name: orders orders_customers_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_customers_id_fk FOREIGN KEY (customers_id) REFERENCES public.customers(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: orders orders_statuses_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_statuses_id_fk FOREIGN KEY (status_id) REFERENCES public.statuses(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: orders orders_transactions_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_transactions_id_fk FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: orders orders_widgets_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_widgets_id_fk FOREIGN KEY (widget_id) REFERENCES public.widgets(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: transactions transactions_transaction_statuses_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: labyrinth
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_transaction_statuses_id_fk FOREIGN KEY (transaction_status_id) REFERENCES public.transaction_statuses(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

