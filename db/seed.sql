INSERT INTO tasks (name, execute, message, hash, active)
VALUES
    ('Sample health check', '0 */5 * * * *', 'Cronjob manager health check', 'seed-sample-health-check', FALSE)
ON CONFLICT (hash) DO NOTHING;
