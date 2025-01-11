create database if not exists webook;
create table if not exists webook.interactives
(
    id          bigint auto_increment
        primary key,
    biz_id      bigint       null,
    biz         varchar(128) null,
    read_cnt    bigint       null,
    collect_cnt bigint       null,
    like_cnt    bigint       null,
    ctime       bigint       null,
    utime       bigint       null,
    constraint biz_type_id
        unique (biz_id, biz)
);

create table if not exists webook.user_collection_bizs
(
    id     bigint auto_increment
        primary key,
    cid    bigint       null,
    biz_id bigint       null,
    biz    varchar(128) null,
    uid    bigint       null,
    ctime  bigint       null,
    utime  bigint       null,
    constraint biz_type_id_uid
        unique (biz_id, biz, uid)
);

create index idx_user_collection_bizs_cid
    on webook.user_collection_bizs (cid);

create table if not exists webook.user_like_bizs
(
    id     bigint auto_increment
        primary key,
    biz_id bigint           null,
    biz    varchar(128)     null,
    uid    bigint           null,
    status tinyint unsigned null,
    ctime  bigint           null,
    utime  bigint           null,
    constraint biz_type_id_uid
        unique (biz_id, biz, uid)
);

