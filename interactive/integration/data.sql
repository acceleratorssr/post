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

INSERT INTO `likes`(`obj_id`, `obj_type`, `view_count`, `collect_count`, `like_count`, `ctime`, `utime`)
VALUES(1,"test",1969,9260,9778,1723257318657,1723257318657),
(2,"test",5597,3779,6794,1723257318657,1723257318657),
(3,"test",6193,3170,5766,1723257318657,1723257318657),
(4,"test",9580,6763,818,1723257318657,1723257318657),
(5,"test",5686,5298,8325,1723257318657,1723257318657),
(6,"test",4387,7005,9872,1723257318657,1723257318657),
(7,"test",5034,4240,1034,1723257318657,1723257318657),
(8,"test",7914,6668,9448,1723257318657,1723257318657),
(9,"test",8217,1252,4860,1723257318657,1723257318657),
(10,"test",5138,1589,1890,1723257318657,1723257318657)