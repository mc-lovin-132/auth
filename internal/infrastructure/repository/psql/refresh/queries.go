package repo

var (
	getQuery = `SELECT value, user_id, device_id, expired_at, revoked, used, session_id 
		FROM refresh_tokens 
		WHERE value = $1;`

	getBySessionIDQuery = `SELECT value, user_id, device_id, expired_at, revoked, used, session_id 
		FROM refresh_tokens 
		WHERE session_id = $1;`

	createQuery = `INSERT INTO refresh_tokens (value, user_id, device_id, expired_at, revoked, used, session_id)
		VALUES (:value, :user_id, :device_id, :expired_at, :revoked, :used, :session_id);`

	markAsUsedQuery = `UPDATE refresh_tokens SET used = TRUE WHERE value = $1;`

	markAsRevokedByDeviceQuery = `UPDATE refresh_tokens SET revoked = TRUE 
		WHERE user_id = $1 AND device_id = $2;`

	markAsRevokedByUserQuery = `UPDATE refresh_tokens SET revoked = TRUE 
		WHERE user_id = $1;`

	markAsRevokedByConcreteQuery = `UPDATE refresh_tokens SET revoked = TRUE 
		WHERE value = $1;`
)
