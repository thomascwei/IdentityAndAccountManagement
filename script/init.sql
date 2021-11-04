CREATE DATABASE IF NOT EXISTS `iam`;
use iam;

CREATE TABLE IF NOT EXISTS `accounts`
(
    `id`       int          NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(30)  NOT NULL,
    `password` VARCHAR(100) NOT NULL,
    `email`    VARCHAR(50)  not NULL,
    `auth`     INT          not NULL,
    UNIQUE (`username`),
    PRIMARY KEY (`id`)
);
# DB初始化時新增管理員帳號 Admin/123456
INSERT accounts
SET username='Admin',
    password='8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',
    Email='admin@admin.com',
    Auth=255;

SELECT *
FROM `accounts`;

# update accounts set password=11111,email='123@abc.com' where id=2;

