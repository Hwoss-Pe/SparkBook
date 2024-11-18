create database webook;
create database webook_intr;
create database webook_article;
create database webook_user;
create database webook_payment;
create database webook_account;
create database webook_reward;
create database webook_comment;
create database webook_tag;

# 准备 canal 用户
CREATE USER 'canal'@'%' IDENTIFIED BY 'canal';
GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' WITH GRANT OPTION;