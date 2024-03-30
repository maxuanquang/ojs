CREATE TABLE IF NOT EXISTS `account` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(50) UNIQUE NOT NULL,
    `role` TINYINT NOT NULL
);

CREATE TABLE IF NOT EXISTS `account_password` (
    `of_account_id` BIGINT UNSIGNED PRIMARY KEY,
    `hashed` VARCHAR(128) NOT NULL,
    FOREIGN KEY (`of_account_id`) REFERENCES `account` (`id`)
);

CREATE TABLE IF NOT EXISTS `problem` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `display_name` VARCHAR(255) NOT NULL,
    `author_id` BIGINT UNSIGNED NOT NULL,
    `description` TEXT NOT NULL,
    `time_limit` INT UNSIGNED NOT NULL,
    `memory_limit` INT UNSIGNED NOT NULL,
    FOREIGN KEY (`author_id`) REFERENCES `account` (`id`)
);

CREATE TABLE IF NOT EXISTS `test_case` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `of_problem_id` BIGINT UNSIGNED NOT NULL,
    `input` TEXT NOT NULL,
    `output` TEXT NOT NULL,
    `is_hidden` TINYINT(1) NOT NULL,
    FOREIGN KEY (`of_problem_id`) REFERENCES `problem` (`id`)
);

CREATE TABLE IF NOT EXISTS `submission` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `of_problem_id` BIGINT UNSIGNED NOT NULL,
    `author_id` BIGINT UNSIGNED NOT NULL,
    `content` TEXT NOT NULL,
    `language` VARCHAR(16) NOT NULL,
    `status` TINYINT NOT NULL,
    `result` TINYINT NOT NULL,
    FOREIGN KEY (`of_problem_id`) REFERENCES `problem` (`id`),
    FOREIGN KEY (`author_id`) REFERENCES `account` (`id`)
);

CREATE TABLE IF NOT EXISTS `token_public_key` (
    `token_public_key_id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `token_public_key_value` VARBINARY(4096) NOT NULL
);