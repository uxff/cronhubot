
# 创建定时任务表
CREATE TABLE `athena_cronjobs` (
                            `id` INT(11) NOT NULL AUTO_INCREMENT,
                            `name` VARCHAR(50) DEFAULT '',
                            `url` VARCHAR(256) NOT NULL DEFAULT '',
                            `expression` VARCHAR(100) NOT NULL DEFAULT '',
                            `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
                            `retries` INT(11) UNSIGNED NOT NULL DEFAULT 0,
                            `request_timeout` INT(11) UNSIGNED NOT NULL DEFAULT 3,
                            `stop_at` TIMESTAMP NOT NULL DEFAULT '1999-01-01 00:00:00',
                            `created_at` TIMESTAMP NOT NULL DEFAULT '1999-01-01 00:00:00',
                            `updated_at` TIMESTAMP NOT NULL DEFAULT '1999-01-01 00:00:00',
                            `group_id` INT(11) UNSIGNED NOT NULL DEFAULT '0',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8;

