SELECT COUNT(id) FROM users WHERE email ILIKE $1 AND id != $2
