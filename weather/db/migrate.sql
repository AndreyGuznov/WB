CREATE TABLE locations (
    id serial unique,
	name VARCHAR(30) not NULL,
	country VARCHAR(30) NOT null,
	lat NUMERIC(7,2) NOT NULL,
	lng NUMERIC(7,3) not null,
	primary key (name, country)
);
CREATE TABLE forecast (
    location_id INT not NULL,
	timestamp BIGINT NOT null,
	temp NUMERIC(7,2) NOT null,
	data text not null,
	primary key (location_id,timestamp),
	FOREIGN KEY (location_id) REFERENCES locations (id) ON DELETE CASCADE
);