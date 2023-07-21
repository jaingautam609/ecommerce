ALTER TABLE users
    ALTER COLUMN user_name DROP NOT NULL,
    alter column user_email drop not null ,
    alter column password drop not null ;
alter table user_role
    alter column user_type drop not null ;
