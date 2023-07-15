create table if not exists users(
    id serial primary key ,
    user_name text not null,
    user_email text not null,
    password BYTEA NOT NULL,
    archive_at timestamp,
    joined_at timestamp default current_timestamp
);
create table if not exists user_role(
    id serial primary key,
    user_id int references users(id),
    user_type text not null
);
create table if not exists item_type(
    id serial primary key,
    item_type text not null,
    added_by int references users(id),
    added_on timestamp default current_timestamp
);
create table if not exists item(
    id serial primary key,
    type_id int references item_type(id),
    item_name text not null,
    added_by int references users(id),
    price int not null,
    archive_at timestamp,
    added_on timestamp default current_timestamp
);
create table if not exists uploads(
    id serial primary key,
    path text not null,
    name text not null,
    url text not null
);
create table if not exists item_image(
    id serial primary key,
    item_id int references item(id),
    upload_id int references uploads(id)
);
create table if not exists cart(
    id serial primary key,
    assign_to int references users(id)
);
create table if not exists cart_item(
    id serial primary key,
    cart_id int references cart(id),
    item_name text not null,
    item_type text not null,
    price int not null ,
    quantity int not null,
    item_id int references item(id),
    added_on timestamp default current_timestamp
);
