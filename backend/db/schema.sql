CREATE TABLE IF NOT EXISTS `side_job` (
    `id`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    `name`	TEXT NOT NULL,
    `minimal_hourly_wage`	INTEGER NOT NULL,
    `maximal_hourly_wage`	INTEGER NOT NULL,
    `minmal_monthly_uptime`	INTEGER NOT NULL,
    `skill_id`	INTEGER NOT NULL,
    `difficulty`	INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS skill (
    `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    `name` TEXT NOT NULL,
    `side_job_id` INTEGER NOT NULL UNIQUE
);