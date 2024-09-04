-- +goose Up
-- +goose StatementBegin
CREATE TABLE `candidates` (
  `candidate_id`      VARCHAR(255) NOT NULL UNIQUE,
  `github_id`         VARCHAR(255) NOT NULL,
  `followers`         INT,
  `contributions`     INT,
  `most_popular_repo` VARCHAR(255),
  `repo_stars`        INT,
  `score`             INT,
  `job_id`            VARCHAR(255),
  PRIMARY KEY (`candidate_id`),
  FOREIGN KEY (`job_id`) REFERENCES `jobs` (`job_id`) ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `candidates`;
-- +goose StatementEnd