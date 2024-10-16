create database if not exists `SSO`;
create database if not exists `garden_user`;
create database if not exists `garden_interactive`;
create database if not exists `garden_article`;
create database if not exists `garden_interactive_partition`;
create database if not exists `garden_gorse`;

create user 'canal'@'%' identified by 'canal';
grant all privileges on *.* to 'canal'@'%' with grant option;
# 只给主从的权限
# GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';

create user 'gorse'@'%' identified by 'gorse_pass';
grant all privileges on gorse.* to 'gorse'@'%';

FLUSH PRIVILEGES;

