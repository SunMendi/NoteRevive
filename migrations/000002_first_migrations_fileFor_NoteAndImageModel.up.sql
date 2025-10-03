create table notes (
     id serial primary key ,
     user_id INTEGER not null, 
     title varchar(255) not null,
     content text not null,
     created_at TIMESTAMP not null DEFAULT NOW(),
     updated_at TIMESTAMP not null DEFAULT NOW()
);

create table note_images(
      id serial primary key,
      note_id INTEGER not null REFERENCES notes(id) on DELETE CASCADE,
      image_url varchar(512) not null,
      uploaded_at TIMESTAMP not null DEFAULT NOW()
);