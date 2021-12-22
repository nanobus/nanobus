CREATE TABLE customers (
	id serial PRIMARY KEY,
	first_name varchar(256) NOT NULL,
	middle_name varchar(256),
	last_name varchar(256) NOT NULL,
	email varchar(256) NOT NULL
);

CREATE TABLE customer_addresses (
	customer_id int NOT NULL REFERENCES customers(id),
	line1 varchar(256) NOT NULL,
	line2 varchar(256),
	city varchar(256) NOT NULL,
	state varchar(2) NOT NULL,
	zip varchar(5) NOT NULL
);