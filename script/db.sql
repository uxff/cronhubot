
# 创建定时任务表
CREATE TABLE `athena_cronjobs` (
                            `id` int(11) NOT NULL AUTO_INCREMENT,
                            `name` varchar(50) DEFAULT '',
                            `url` varchar(256) NOT NULL DEFAULT '',
                            `expression` varchar(100) NOT NULL DEFAULT '',
                            `status` tinyint DEFAULT 1,
                            `retries` int(11) DEFAULT '0',
                            `request_timeout` int(11) DEFAULT '3',
                            `stop_at` timestamp NOT NULL DEFAULT '1999-01-01 00:00:00',
                            `created_at` timestamp NOT NULL DEFAULT '1999-01-01 00:00:00',
                            `updated_at` timestamp NOT NULL DEFAULT '1999-01-01 00:00:00',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8;

