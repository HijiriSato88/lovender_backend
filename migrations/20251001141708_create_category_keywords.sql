-- Create "category_keywords" table
CREATE TABLE `category_keywords` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` smallint unsigned NOT NULL,
  `keyword` varchar(100) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  INDEX `idx_category_keywords_category` (`category_id`),
  UNIQUE INDEX `uq_category_keywords` (`category_id`, `keyword`),
  CONSTRAINT `fk_category_keywords_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
