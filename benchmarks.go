package main

import (
	"database/sql"
)

func bmSimpleExec(db *sql.DB, n int) error {
	for i := 0; i < n; i++ {
		if _, err := db.Exec("DO 1"); err != nil {
			return err
		}
	}
	return nil
}

func bmPreparedExec(db *sql.DB, n int) error {
	stmt, err := db.Prepare("DO 1")
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if _, err := stmt.Exec(); err != nil {
			return err
		}
	}

	return stmt.Close()
}

func bmSimpleQueryRow(db *sql.DB, n int) error {
	var num int

	for i := 0; i < n; i++ {
		if err := db.QueryRow("SELECT 1").Scan(&num); err != nil {
			return err
		}
	}
	return nil
}

func bmPreparedQueryRow(db *sql.DB, n int) error {
	var num int

	stmt, err := db.Prepare("SELECT 1")
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if err := stmt.QueryRow().Scan(&num); err != nil {
			return err
		}
	}

	return stmt.Close()
}

func bmPreparedQueryRowParam(db *sql.DB, n int) error {
	var num int

	stmt, err := db.Prepare("SELECT ?")
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if err := stmt.QueryRow(i).Scan(&num); err != nil {
			return err
		}
	}

	return stmt.Close()
}

func bmEchoMixed5(db *sql.DB, n int) error {
	stmt, err := db.Prepare("SELECT ?, ?, ?, ?, ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Some random data with different types
	type entry struct {
		id    int64
		name  string
		ratio float64
		other interface{}
		hire  bool
	}

	in := entry{
		id:    42,
		name:  "Gopher",
		ratio: 1.618,
		other: nil,
		hire:  true,
	}

	var out entry

	for i := 0; i < n; i++ {
		if err := stmt.QueryRow(
			in.id,
			in.name,
			in.ratio,
			in.other,
			in.hire,
		).Scan(
			&out.id,
			&out.name,
			&out.ratio,
			&out.other,
			&out.hire,
		); err != nil {
			return err
		}
	}
	return nil
}

func bmSelectLargeString(db *sql.DB, n int) error {
	var str string
	for i := 0; i < n; i++ {
		if err := db.QueryRow("SELECT REPEAT('A', 10000)").Scan(&str); err != nil {
			return err
		}
	}
	return nil
}

func bmSelectPreparedLargeString(db *sql.DB, n int) error {
	stmt, err := db.Prepare("SELECT REPEAT('A', 10000)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var str string
	for i := 0; i < n; i++ {
		if err := stmt.QueryRow().Scan(&str); err != nil {
			return err
		}
	}
	return nil
}

func bmSelectLargeBytes(db *sql.DB, n int) error {
	var raw []byte
	for i := 0; i < n; i++ {
		if err := db.QueryRow("SELECT REPEAT('A', 10000)").Scan(&raw); err != nil {
			return err
		}
	}
	return nil
}

func bmSelectPreparedLargeBytes(db *sql.DB, n int) error {
	stmt, err := db.Prepare("SELECT REPEAT('A', 10000)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var raw []byte
	for i := 0; i < n; i++ {
		if err := stmt.QueryRow().Scan(&raw); err != nil {
			return err
		}
	}
	return nil
}

func bmSelectLargeRaw(db *sql.DB, n int) error {
	var raw sql.RawBytes
	for i := 0; i < n; i++ {
		rows, err := db.Query("SELECT REPEAT('A', 10000)")
		if err != nil {
			return err
		}

		if !rows.Next() {
			return sql.ErrNoRows
		}

		if err = rows.Scan(&raw); err != nil {
			return err
		}

		if err = rows.Close(); err != nil {
			return err
		}
	}
	return nil
}

func bmSelectPreparedLargeRaw(db *sql.DB, n int) error {
	stmt, err := db.Prepare("SELECT REPEAT('A', 10000)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var raw sql.RawBytes
	for i := 0; i < n; i++ {
		rows, err := stmt.Query()
		if err != nil {
			return err
		}

		if !rows.Next() {
			return sql.ErrNoRows
		}

		if err = rows.Scan(&raw); err != nil {
			return err
		}

		if err = rows.Close(); err != nil {
			return err
		}
	}
	return nil
}