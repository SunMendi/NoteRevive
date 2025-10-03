create table users (
     id serial primary key,
     name varchar(255) not null,
     google_id varchar(255) not null,
     email varchar(255) not null,

     birthday date,
     gender varchar(50),
     created_at TIMESTAMP not null DEFAULT NOW(),
     updated_at TIMESTAMP not null DEFAULT NOW(),
     deleted_at TIMESTAMP

);

create UNIQUE index idx_users_email on users(email);