-- +goose Up
-- +goose StatementBegin
CREATE TABLE `jobs` (
  `job_id`        VARCHAR(255) NOT NULL UNIQUE,
  `job_name`      VARCHAR(255) NOT NULL,
  `description`   VARCHAR(255) NOT NULL,
  `drive_link`    VARCHAR(255),
  `cloud_provider` VARCHAR(255),
  `created_at`    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `user_id`       VARCHAR(255),
  PRIMARY KEY (`job_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `jobs`;
-- +goose StatementEnd
