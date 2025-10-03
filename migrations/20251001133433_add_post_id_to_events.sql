-- Modify "events" table
ALTER TABLE `events` ADD COLUMN `post_id` bigint unsigned NULL AFTER `category_id`;
