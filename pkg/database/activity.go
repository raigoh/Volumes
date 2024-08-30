package database

import (
	"fmt"
	"literary-lions-forum/internal/models"
	"time"
)

// GetActiveUsers returns the number of users who have been active in the last n days
func GetActiveUsers(days int) (int, error) {
	var count int
	// Calculate the cutoff date for activity
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// SQL query to count distinct users who have either posted or commented since the cutoff date
	query := `
		SELECT COUNT(DISTINCT user_id) 
		FROM (
			SELECT user_id FROM posts WHERE created_at > ?
			UNION
			SELECT user_id FROM comments WHERE created_at > ?
		)
	`
	// Execute the query and scan the result into count
	err := DB.QueryRow(query, cutoffDate, cutoffDate).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting active users: %v", err)
	}
	return count, nil
}

// GetRecentActivity returns the n most recent activities (posts or comments)
func GetRecentActivity(n int) ([]models.Activity, error) {
	// SQL query to fetch recent posts and comments, ordered by creation date
	query := `
		SELECT 'post' as type, u.username, p.title as content, p.created_at as created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		UNION ALL
		SELECT 'comment' as type, u.username, c.content, c.created_at as created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		ORDER BY created_at DESC
		LIMIT ?
	`
	// Execute the query
	rows, err := DB.Query(query, n)
	if err != nil {
		return nil, fmt.Errorf("error querying recent activity: %v", err)
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var a models.Activity
		var createdAt time.Time
		// Scan each row into an Activity struct
		err := rows.Scan(&a.Type, &a.Username, &a.Content, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning activity row: %v", err)
		}
		// Format the creation date
		a.Date = createdAt.Format("02.01.2006 15:04")
		activities = append(activities, a)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activity rows: %v", err)
	}

	return activities, nil
}
