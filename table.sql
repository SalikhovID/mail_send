CREATE TABLE mail_send_message (
    id serial PRIMARY KEY,
    company_id integer DEFAULT 0 NOT NULL,
    email character varying(100) NOT NULL,
    create_time timestamp without time zone,
    send_time timestamp without time zone,
    mess_title character varying(250) NOT NULL,
    mess_body text NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    from_col character varying(100)
);