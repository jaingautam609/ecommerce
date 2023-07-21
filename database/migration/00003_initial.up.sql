alter table users
    add column phone_no text ,
    add column is_verified_by_phone bool,
    add column is_verified_by_email bool;
create table if not exists opt(
    id serial primary key ,
    user_id int references users(id),
    phone_no text ,
    user_email text,
    otp text ,
    created_at timestamp default current_timestamp,
    expired_at timestamp default current_timestamp + INTERVAL '10 minutes'
)