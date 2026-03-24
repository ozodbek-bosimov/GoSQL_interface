package models

import (
	"database/sql"
	"fmt"
	"time"
)

// ACL represents an access control list rule
type ACL struct {
	ID        int64     `json:"id"`
	SrcIP     string    `json:"src_ip"`
	DstIP     string    `json:"dst_ip"`
	Protocol  string    `json:"protocol"`
	SrcPort   *string   `json:"src_port"`
	DstPort   *string   `json:"dst_port"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

// GetAllACLs retrieves all ACL rules with optional search
func GetAllACLs(db *sql.DB, search string) ([]ACL, error) {
	query := `
		SELECT id, src_ip, dst_ip, protocol, src_port, dst_port, action, created_at
		FROM acl
		WHERE ($1 = '' OR src_ip ILIKE $2 OR dst_ip ILIKE $2)
		ORDER BY id DESC
	`

	searchPattern := "%" + search + "%"
	rows, err := db.Query(query, search, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to query ACL rules: %w", err)
	}
	defer rows.Close()

	var acls []ACL
	for rows.Next() {
		var a ACL
		err := rows.Scan(&a.ID, &a.SrcIP, &a.DstIP, &a.Protocol, &a.SrcPort, &a.DstPort, &a.Action, &a.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ACL: %w", err)
		}
		acls = append(acls, a)
	}

	return acls, nil
}

// GetACLByID retrieves an ACL rule by ID
func GetACLByID(db *sql.DB, id int64) (*ACL, error) {
	query := `SELECT id, src_ip, dst_ip, protocol, src_port, dst_port, action, created_at FROM acl WHERE id = $1`

	var a ACL
	err := db.QueryRow(query, id).Scan(&a.ID, &a.SrcIP, &a.DstIP, &a.Protocol, &a.SrcPort, &a.DstPort, &a.Action, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ACL rule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ACL rule: %w", err)
	}

	return &a, nil
}

// CreateACL creates a new ACL rule
func CreateACL(db *sql.DB, acl *ACL) error {
	query := `
		INSERT INTO acl (src_ip, dst_ip, protocol, src_port, dst_port, action)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := db.QueryRow(query, acl.SrcIP, acl.DstIP, acl.Protocol, acl.SrcPort, acl.DstPort, acl.Action).
		Scan(&acl.ID, &acl.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create ACL rule: %w", err)
	}

	return nil
}

// UpdateACL updates an existing ACL rule
func UpdateACL(db *sql.DB, acl *ACL) error {
	query := `
		UPDATE acl
		SET src_ip = $1, dst_ip = $2, protocol = $3, src_port = $4, dst_port = $5, action = $6
		WHERE id = $7
	`

	_, err := db.Exec(query, acl.SrcIP, acl.DstIP, acl.Protocol, acl.SrcPort, acl.DstPort, acl.Action, acl.ID)
	if err != nil {
		return fmt.Errorf("failed to update ACL rule: %w", err)
	}

	return nil
}

// DeleteACL deletes an ACL rule by ID
func DeleteACL(db *sql.DB, id int64) error {
	query := `DELETE FROM acl WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ACL rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ACL rule not found")
	}

	return nil
}
