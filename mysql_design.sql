
create table if not exists users (
    id int not null unique auto_increment,
    username varchar (127) not null unique,
    password varchar (127) not null,
    primary key (id)
);

create table if not exists links (
    id int not null unique auto_increment,
    title varchar (255) ,
    address varchar (255) ,
    user_id int ,
    foreign key (user_id) references users(id),
    primary key (id)
);

create table if not exists projects (
    id int not null unique auto_increment,
    name varchar (255) not null,
    quantity int,
    buget decimal(10, 2),
    primary key (id)
);

create table if not exists project_users (
    project_id int not null,
    user_id int not null,
    foreign key (project_id) references projects(id) on delete cascade,
    foreign key (user_id) references users(id) on delete cascade,
    primary key (project_id, user_id)
);