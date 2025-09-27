package repository

import (
	"context"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TopUpRepository struct {
	db *pgxpool.Pool
}

func NewTopUpRepository(db *pgxpool.Pool) *TopUpRepository {
	return &TopUpRepository{db: db}
}

func (tr *TopUpRepository) CreateTopUp(ctx context.Context, topup *models.TopUp) (*models.TopUp, error) {
	query := `
		INSERT INTO topup (amount, tax, payment_id, topup_status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	err := tr.db.QueryRow(ctx, query,
		topup.Amount,
		topup.Tax,
		topup.PaymentID,
		topup.Status,
	).Scan(&topup.ID, &topup.CreatedAt)
	if err != nil {
		return nil, err
	}
	return topup, nil
}

func (tr *TopUpRepository) UpdateStatusTopUp(c context.Context, topupID int, status models.TopUpStatus) error {
	query := `UPDATE topup SET topup_status = $1, updated_at = NOW() WHERE id = $2`
	_, err := tr.db.Exec(c, query, status, topupID)
	return err
}

func (tr *TopUpRepository) GetTopUpByID(c context.Context, topupID int) (*models.TopUp, error) {
	query := `SELECT id, amount, tax, payment_id, topup_status, created_at, updated_at FROM topup WHERE id = $1`
	row := tr.db.QueryRow(c, query, topupID)

	var t models.TopUp
	err := row.Scan(&t.ID, &t.Amount, &t.Tax, &t.PaymentID, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (tr *TopUpRepository) ApplyToWallet(c context.Context, walletID int, topupID int, amount int) error {
	tx, err := tr.db.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	qInsertIntoWalletTopup := "INSERT INTO wallets_topup (wallets_id, topup_id) VALUES ($1, $2)"
	_, err = tx.Exec(c, qInsertIntoWalletTopup, walletID, topupID)
	if err != nil {
		return err
	}

	qUpdateToWallets := "UPDATE wallets SET balance = balance + $1 WHERE id = $2"
	_, err = tx.Exec(c, qUpdateToWallets, amount, walletID)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func (tr *TopUpRepository) GetWalletIDByUserID(c context.Context, userID int) (int, error) {
	query := `SELECT id FROM wallets WHERE user_id = $1 LIMIT 1`
	row := tr.db.QueryRow(c, query, userID)

	var walletID int
	if err := row.Scan(&walletID); err != nil {
		return 0, err
	}
	return walletID, nil
}

func (tr *TopUpRepository) FindAllPaymentMethods(c context.Context) ([]models.PaymentMethod, error) {
	query := `SELECT id, name FROM payment_method ORDER BY id ASC`
	rows, err := tr.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []models.PaymentMethod
	for rows.Next() {
		var pm models.PaymentMethod
		if err := rows.Scan(&pm.ID, &pm.Name); err != nil {
			return nil, err
		}
		methods = append(methods, pm)
	}
	return methods, nil
}

// ini yang utama
func (tr *TopUpRepository) CreateTopUpTransaction(ctx context.Context, topup *models.TopUp, userID int) (*models.TopUp, error) {
	tx, err := tr.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	queryInsertTopup := `
		INSERT INTO topup (amount, tax, payment_id, topup_status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, queryInsertTopup,
		topup.Amount,
		topup.Tax,
		topup.PaymentID,
		models.TopUpSuccess,
	).Scan(&topup.ID, &topup.CreatedAt)
	if err != nil {
		return nil, err
	}

	var walletID int
	queryWallet := `SELECT id FROM wallets WHERE user_id = $1 LIMIT 1`
	if err := tx.QueryRow(ctx, queryWallet, userID).Scan(&walletID); err != nil {
		return nil, err
	}

	qInsertWalletTopup := `INSERT INTO wallets_topup (wallets_id, topup_id) VALUES ($1, $2)`
	if _, err := tx.Exec(ctx, qInsertWalletTopup, walletID, topup.ID); err != nil {
		return nil, err
	}

	qUpdateWallet := `UPDATE wallets SET balance = balance + $1 WHERE id = $2`
	if _, err := tx.Exec(ctx, qUpdateWallet, topup.Amount, walletID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	topup.Status = models.TopUpSuccess
	return topup, nil
}
