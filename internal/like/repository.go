package like

import (
	"database/sql"
	"fmt"
	"literary-lions-forum/pkg/database"
)

// AddLike adds or updates a like/dislike
func AddLike(userID int, targetID int, targetType string, isLike bool) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Check if a like already exists
	var existingID int
	var existingIsLike bool
	query := `
			SELECT id, is_like FROM likes 
			WHERE user_id = ? AND 
			((?='post' AND post_id=?) OR (?='comment' AND comment_id=?))
	`
	err = tx.QueryRow(query, userID, targetType, targetID, targetType, targetID).Scan(&existingID, &existingIsLike)

	if err == sql.ErrNoRows {
		// No existing like, insert a new one
		insertQuery := `
					INSERT INTO likes (user_id, post_id, comment_id, is_like)
					VALUES (?, 
							CASE WHEN ? = 'post' THEN ? ELSE NULL END, 
							CASE WHEN ? = 'comment' THEN ? ELSE NULL END, 
							?
					)
			`
		_, err = tx.Exec(insertQuery, userID, targetType, targetID, targetType, targetID, isLike)
	} else if err == nil {
		// Existing like found, update it
		updateQuery := `UPDATE likes SET is_like = ? WHERE id = ?`
		_, err = tx.Exec(updateQuery, isLike, existingID)
	} else {
		return fmt.Errorf("error checking for existing like: %v", err)
	}

	if err != nil {
		return fmt.Errorf("error adding/updating like: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// RemoveLike removes a like/dislike
func RemoveLike(userID int, targetID int, targetType string) error {
	query := `
			DELETE FROM likes
			WHERE user_id = ? AND
			CASE WHEN ? = 'post' THEN post_id = ? ELSE comment_id = ? END
	`
	_, err := database.DB.Exec(query, userID, targetType, targetID, targetID)
	return err
}

// GetLikesCount returns the number of likes and dislikes for a target
func GetLikesCount(targetID int, targetType string) (likes int, dislikes int, err error) {
	query := `
			SELECT 
					SUM(CASE WHEN is_like = 1 THEN 1 ELSE 0 END) as likes,
					SUM(CASE WHEN is_like = 0 THEN 1 ELSE 0 END) as dislikes
			FROM likes
			WHERE (? = 'post' AND post_id = ?) OR (? = 'comment' AND comment_id = ?)
	`
	err = database.DB.QueryRow(query, targetType, targetID, targetType, targetID).Scan(&likes, &dislikes)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting like counts: %v", err)
	}
	return
}

// GetUserLike returns the user's like status for a target
func GetUserLike(userID int, targetID int, targetType string) (isLike *bool, err error) {
	query := `
			SELECT is_like
			FROM likes
			WHERE user_id = ? AND
			CASE WHEN ? = 'post' THEN post_id = ? ELSE comment_id = ? END
	`
	var like bool
	err = database.DB.QueryRow(query, userID, targetType, targetID, targetID).Scan(&like)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &like, nil
}
