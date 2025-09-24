-- Create "categories" table
CREATE TABLE `categories` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT,
  `slug` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` text NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  INDEX `idx_categories_slug` (`slug`),
  UNIQUE INDEX `uq_categories_slug` (`slug`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "users" table
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uq_users_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "oshis" table
CREATE TABLE `oshis` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `name` varchar(191) NOT NULL,
  `description` text NULL,
  `theme_color` char(7) NOT NULL DEFAULT "#FFFFFF",
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  INDEX `idx_oshis_user` (`user_id`),
  UNIQUE INDEX `uq_oshis_user_name` (`user_id`, `name`),
  CONSTRAINT `fk_oshis_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "events" table
CREATE TABLE `events` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `oshi_id` bigint unsigned NOT NULL,
  `category_id` smallint unsigned NULL,
  `title` varchar(255) NOT NULL,
  `description` text NULL,
  `has_alarm` bool NOT NULL DEFAULT 1,
  `notification_timing` enum('0','5m','10m','15m','30m','1h','2h','1d','2d','1w') NULL DEFAULT "15m",
  `has_notification_sent` bool NULL DEFAULT 0,
  `starts_at` datetime(3) NOT NULL,
  `ends_at` datetime(3) NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  INDEX `idx_events_category` (`category_id`),
  INDEX `idx_events_oshi` (`oshi_id`),
  INDEX `idx_events_oshi_starts` (`oshi_id`, `starts_at`),
  CONSTRAINT `fk_events_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT `fk_events_oshi` FOREIGN KEY (`oshi_id`) REFERENCES `oshis` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `chk_events_time` CHECK ((`ends_at` is null) or (`ends_at` >= `starts_at`))
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "oshi_accounts" table
CREATE TABLE `oshi_accounts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `oshi_id` bigint unsigned NOT NULL,
  `url` varchar(2048) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  INDEX `idx_oshi_accounts_oshi` (`oshi_id`),
  UNIQUE INDEX `uq_oshi_accounts` (`oshi_id`, `url` (191)),
  CONSTRAINT `fk_oshi_accounts_oshi` FOREIGN KEY (`oshi_id`) REFERENCES `oshis` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "oshi_categories" table
CREATE TABLE `oshi_categories` (
  `oshi_id` bigint unsigned NOT NULL,
  `category_id` smallint unsigned NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`oshi_id`, `category_id`),
  INDEX `fk_oc_category` (`category_id`),
  INDEX `idx_oshi_categories_oshi` (`oshi_id`),
  CONSTRAINT `fk_oc_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_oc_oshi` FOREIGN KEY (`oshi_id`) REFERENCES `oshis` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
