INSERT INTO users (user_name, user_email, password)
VALUES ('Gautam Jain', 'gautamjain@example.com', '$2a$04$L4AUohH6PwpDikwqh8w68eJ/v566vc9RFOejyQNwMiSMfXoNxx6la');

WITH users_id AS (
    SELECT id
    FROM users
    WHERE user_name = 'Gautam Jain'
    LIMIT 1
)
INSERT INTO user_role (user_id, user_type)
SELECT id, 'admin'
FROM users_id;
