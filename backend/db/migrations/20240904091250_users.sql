-- +goose Up
-- +goose StatementBegin
CREATE TABLE `users` (
  `user_id`         VARCHAR(255) NOT NULL UNIQUE,
  `first_name`      VARCHAR(255) NOT NULL,
  `last_name`       VARCHAR(255) NOT NULL,
  `username`        VARCHAR(255) NOT NULL UNIQUE,
  `password`        VARCHAR(255) NOT NULL,
  `email`           VARCHAR(255) NOT NULL UNIQUE,
  `organization`    VARCHAR(255) NOT NULL,
  `created_at`      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB;;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `users`;
-- +goose StatementEnd
