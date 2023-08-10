package db

import (
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func PrepareDB(ctx context.Context) (*sql.DB, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current path: %w")
	}

	db, err := sql.Open("sqlite3", filepath.Join(path, "db", "sideJob.sqlite3"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DB: %w")
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping DB: %w")
	}

	f, err := os.ReadFile(filepath.Join(path, "db", "schema.sql"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open schema.sql %w")
	}

	if _, err = db.ExecContext(ctx, string(f)); err != nil {
		return nil, errors.Wrap(err, "failed to exec query: %w")
	}

	tx, err := db.Begin();if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction: %w")
    }
	defer tx.Rollback()

	// プリペアドステートメントを作成する
    stmt, err := tx.Prepare("INSERT INTO side_job VALUES (?, ?, ?)"); if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    // CSVファイルを開く
    file, err := os.Open("./side_job.csv"); if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // CSVファイルを読み込む
    reader := csv.NewReader(file)

    // ヘッダー行を読み飛ばす
    _, err = reader.Read(); if err != nil {
        log.Fatal(err)
    }

	// データ行を読み込んでデータベースに挿入する
	for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

		// CSVの各列をinterface{}型のスライスに変換する
		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = v
		}

		// プリペアドステートメントに値をバインドして実行する
		_, err = stmt.Exec(args...)
		if err != nil {
			log.Fatal(err)
		}
	}

	return db, nil
}
