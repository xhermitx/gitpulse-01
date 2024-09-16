-- +goose Up
-- +goose StatementBegin
CREATE TABLE `candidates` (
  `candidate_id`                           VARCHAR(255) NOT NULL UNIQUE,
  `name`                                   VARCHAR(255) NOT NULL,
  `username`                               VARCHAR(255) NOT NULL,
  `avatar_url`                             VARCHAR(255),
  `bio`                                    VARCHAR(255),
  `email`                                  VARCHAR(255),
  `website_url`                            VARCHAR(255),
  `total_contributions`                    INT,
  `total_followers`                        INT,
  `top_repo`                               VARCHAR(255),
  `top_repo_stars`                         INT,
  `top_contributed_repo`                   VARCHAR(255),
  `top_contributed_repo_stars`             INT,
  `languages`                              INT,
  `topics`                                 INT,
  `job_id`                                 VARCHAR(255),
  PRIMARY KEY (`candidate_id`),
  FOREIGN KEY (`job_id`) REFERENCES `jobs` (`job_id`) ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `candidates`;
-- +goose StatementEnd