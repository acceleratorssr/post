create database if not exists test;
create table if not exists test.likes
(
    id              bigint auto_increment primary key,
    obj_id          bigint       null,
    obj_type        varchar(128) null,
    view_count      bigint       null,
    collect_count   bigint       null,
    like_count      bigint       null,
    ctime           bigint       null,
    utime           bigint       null,
    constraint obj_id_type
        unique (obj_id, obj_type)
);

create table if not exists test.user_give_collects
(
    id          bigint auto_increment primary key,
    obj_id      bigint       null,
    obj_type    varchar(128) null,
    uid         bigint       null,
    ctime       bigint       null,
    utime       bigint       null,
    constraint obj_uid_id_type
        unique (uid, obj_id, obj_type)
);

create table if not exists test.user_give_likes
(
    id          bigint auto_increment primary key,
    obj_id      bigint           null,
    obj_type    varchar(128)     null,
    uid         bigint           null,
    status      tinyint unsigned null,
    ctime       bigint           null,
    utime       bigint           null,
    constraint obj_uid_id_type
        unique (uid, obj_id, obj_type)
);

