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
