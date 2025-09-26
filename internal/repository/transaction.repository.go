// repository/transaction.go (Updated with Soft Delete)
package repository

import (
	"context"
	"errors"
	"log"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (tr *TransactionRepository) GetHistory(ctx context.Context, userID int) ([]models.TransactionHistory, error) {
	sql := `SELECT 
		t.id,
		CASE 
			WHEN w_receiver.user_id = $1 THEN 'Transfer'
			ELSE 'Send'
		END as transaction_type,

		CASE 
			WHEN w_receiver.user_id = $1 THEN
				COALESCE(p_sender.profile_picture, '')  
			ELSE
				COALESCE(p_receiver.profile_picture, '')  
		END as profile_picture,

		CASE 
			WHEN w_receiver.user_id = $1 THEN
				COALESCE(p_sender.fullname, 'Unknown')  
			ELSE
				COALESCE(p_receiver.fullname, 'Unknown')  
		END as contact_name,

		CASE 
			WHEN w_receiver.user_id = $1 THEN
				COALESCE(p_sender.phone, 'Unknown')  
			ELSE
				COALESCE(p_receiver.phone, 'Unknown')  
		END as phone_number,

		CASE 
			WHEN w_receiver.user_id = $1 THEN
				CONCAT('Rp ', TO_CHAR(t.amount, 'FM999,999,999')) 
			ELSE
				CONCAT('Rp ', TO_CHAR(t.amount, 'FM999,999,999'))
		END as display_amount,

		t.amount as original_amount,
		COALESCE(t.transfer_status::text, 'pending') as status,
		COALESCE(t.notes, '') as notes,
		t.created_at

	FROM transfer t
	JOIN wallets w_sender ON t.sender_wallet_id = w_sender.id
	JOIN wallets w_receiver ON t.receiver_wallet_id = w_receiver.id
	LEFT JOIN profile p_sender ON w_sender.user_id = p_sender.user_id
	LEFT JOIN profile p_receiver ON w_receiver.user_id = p_receiver.user_id
	WHERE 
		(w_sender.user_id = $1 OR w_receiver.user_id = $1)
		AND (
			CASE 
				WHEN w_sender.user_id = $1 THEN t.deleted_by_sender = FALSE
				WHEN w_receiver.user_id = $1 THEN t.deleted_by_receiver = FALSE
			END
		)
	ORDER BY t.created_at DESC`

	rows, err := tr.db.Query(ctx, sql, userID)
	if err != nil {
		log.Printf("Error querying transaction history: %v", err)
		return nil, err
	}
	defer rows.Close()

	var histories []models.TransactionHistory
	for rows.Next() {
		var history models.TransactionHistory
		if err := rows.Scan(
			&history.ID,
			&history.Type,
			&history.ProfilePicture,
			&history.ContactName,
			&history.PhoneNumber,
			&history.Amount,
			&history.OriginalAmount,
			&history.Status,
			&history.Notes,
			&history.CreatedAt,
		); err != nil {
			log.Printf("Error scanning transaction row: %v", err)
			return nil, err
		}
		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating transaction rows: %v", err)
		return nil, err
	}

	if len(histories) == 0 {
		return nil, errors.New("no transactions found")
	}

	return histories, nil
}

// SoftDeleteTransaction - untuk soft delete transaksi
func (tr *TransactionRepository) SoftDeleteTransaction(ctx context.Context, transactionID, userID int) error {
	// Cek apakah user adalah sender atau receiver
	checkSQL := `SELECT 
		CASE 
			WHEN w_sender.user_id = $2 THEN 'sender'
			WHEN w_receiver.user_id = $2 THEN 'receiver'
			ELSE 'none'
		END as user_role
	FROM transfer t
	JOIN wallets w_sender ON t.sender_wallet_id = w_sender.id
	JOIN wallets w_receiver ON t.receiver_wallet_id = w_receiver.id
	WHERE t.id = $1`

	var userRole string
	err := tr.db.QueryRow(ctx, checkSQL, transactionID, userID).Scan(&userRole)
	if err != nil {
		log.Printf("Error checking user role: %v", err)
		return err
	}

	if userRole == "none" {
		return errors.New("transaction not found or user not authorized")
	}

	// Update berdasarkan role user
	var updateSQL string
	if userRole == "sender" {
		updateSQL = `UPDATE transfer SET deleted_by_sender = TRUE WHERE id = $1`
	} else {
		updateSQL = `UPDATE transfer SET deleted_by_receiver = TRUE WHERE id = $1`
	}

	_, err = tr.db.Exec(ctx, updateSQL, transactionID)
	if err != nil {
		log.Printf("Error soft deleting transaction: %v", err)
		return err
	}

	return nil
}

// GetTopupHistory - untuk mendapatkan history topup
func (tr *TransactionRepository) GetTopupHistory(ctx context.Context, userID int) ([]models.TransactionHistory, error) {
	sql := `SELECT 
		t.id,
		'Topup' as transaction_type,
		'' as profile_picture,
		pm.name as contact_name,
		'' as phone_number,
		CONCAT('+Rp ', TO_CHAR(t.amount, 'FM999,999,999')) as display_amount,
		t.amount as original_amount,
		COALESCE(t.topup_status::text, 'pending') as status,
		CONCAT('Tax: Rp ', TO_CHAR(COALESCE(t.tax, 0), 'FM999,999,999')) as notes,
		t.created_at
	FROM topup t
	JOIN wallets_topup wt ON t.id = wt.topup_id
	JOIN wallets w ON wt.wallets_id = w.id
	JOIN payment_method pm ON t.payment_id = pm.id
	WHERE w.user_id = $1
	ORDER BY t.created_at DESC`

	rows, err := tr.db.Query(ctx, sql, userID)
	if err != nil {
		log.Printf("Error querying topup history: %v", err)
		return nil, err
	}
	defer rows.Close()

	var histories []models.TransactionHistory
	for rows.Next() {
		var history models.TransactionHistory
		if err := rows.Scan(
			&history.ID,
			&history.Type,
			&history.ProfilePicture,
			&history.ContactName,
			&history.PhoneNumber,
			&history.Amount,
			&history.OriginalAmount,
			&history.Status,
			&history.Notes,
			&history.CreatedAt,
		); err != nil {
			log.Printf("Error scanning topup row: %v", err)
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, nil
}

// GetAllHistory - untuk mendapatkan gabungan transfer dan topup history
func (tr *TransactionRepository) GetAllHistory(ctx context.Context, userID int) ([]models.TransactionHistory, error) {
	// Get transfer history
	transferHistory, err := tr.GetHistory(ctx, userID)
	if err != nil && err.Error() != "no transactions found" {
		return nil, err
	}

	// Get topup history
	topupHistory, err := tr.GetTopupHistory(ctx, userID)
	if err != nil {
		log.Printf("Error getting topup history: %v", err)
		// Continue even if topup history fails
		topupHistory = []models.TransactionHistory{}
	}

	// Combine histories
	allHistory := append(transferHistory, topupHistory...)

	// Sort by created_at DESC
	if len(allHistory) > 1 {
		for i := 0; i < len(allHistory)-1; i++ {
			for j := i + 1; j < len(allHistory); j++ {
				if allHistory[i].CreatedAt.Before(allHistory[j].CreatedAt) {
					allHistory[i], allHistory[j] = allHistory[j], allHistory[i]
				}
			}
		}
	}

	if len(allHistory) == 0 {
		return nil, errors.New("no transactions found")
	}

	return allHistory, nil
}
