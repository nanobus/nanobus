CREATE TABLE customers (
	id serial PRIMARY KEY,
	first_name varchar(256) NOT NULL,
	middle_name varchar(256),
	last_name varchar(256) NOT NULL,
	email varchar(256) NOT NULL
);
