-- Modify "events" table
ALTER TABLE `events` ADD COLUMN `url` varchar(2048) NULL AFTER `description`;
