CREATE TABLE temp_offers (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_id varchar(100) UNIQUE ,
    title varchar(1000),
    content text,
    created_at timestamp DEFAULT now()
);

CREATE TABLE temp_companies (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(1000),
    description text
);
