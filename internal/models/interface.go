package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Interface struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IP        string    `json:"ip"`
	MAC       string    `json:"mac"`
	MTU       int       `json:"mtu"`
	Status    bool      `json:"status"`
	IPType    bool      `json:"ip_type"`
	CreatedAt time.Time `json:"created_at"`
}

func GetAllInterfaces(db *sql.DB, search string) ([]Interface, error) {
	query := `
		SELECT id, name, ip, mac, mtu, status, ip_type, created_at
		FROM interfaces
		WHERE ($1 = '' OR name ILIKE $2 OR ip ILIKE $2)
		ORDER BY id DESC
	`

	searchPattern := "%" + search + "%"
	rows, err := db.Query(query, search, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to query interfaces: %w", err)
	}
	defer rows.Close()

	var interfaces []Interface
	for rows.Next() {
		var i Interface
		err := rows.Scan(&i.ID, &i.Name, &i.IP, &i.MAC, &i.MTU, &i.Status, &i.IPType, &i.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interface: %w", err)
		}
		interfaces = append(interfaces, i)
	}

	return interfaces, nil
}

func GetInterfaceByID(db *sql.DB, id int64) (*Interface, error) {
	query := `SELECT id, name, ip, mac, mtu, status, ip_type, created_at FROM interfaces WHERE id = $1`

	var i Interface
	err := db.QueryRow(query, id).Scan(&i.ID, &i.Name, &i.IP, &i.MAC, &i.MTU, &i.Status, &i.IPType, &i.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("interface not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get interface: %w", err)
	}

	return &i, nil
}

func CreateInterface(db *sql.DB, iface *Interface) error {
	query := `
		INSERT INTO interfaces (name, ip, mac, mtu, status, ip_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := db.QueryRow(query, iface.Name, iface.IP, iface.MAC, iface.MTU, iface.Status, iface.IPType).
		Scan(&iface.ID, &iface.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create interface: %w", err)
	}

	return nil
}

func UpdateInterface(db *sql.DB, iface *Interface) error {
	query := `
		UPDATE interfaces
		SET name = $1, ip = $2, mac = $3, mtu = $4, status = $5, ip_type = $6
		WHERE id = $7
	`

	_, err := db.Exec(query, iface.Name, iface.IP, iface.MAC, iface.MTU, iface.Status, iface.IPType, iface.ID)
	if err != nil {
		return fmt.Errorf("failed to update interface: %w", err)
	}

	return nil
}

func DeleteInterface(db *sql.DB, id int64) error {
	query := `DELETE FROM interfaces WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete interface: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("interface not found")
	}

	return nil
}
