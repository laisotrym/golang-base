CREATE TABLE IF NOT EXISTS `state_sys` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `last_time` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;