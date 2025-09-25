// repository/transaction.go (Updated)
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
				CONCAT('+Rp', TO_CHAR(t.amount, 'FM999,999')) 
			ELSE
				CONCAT('-Rp', TO_CHAR(t.amount, 'FM999,999'))
		END as display_amount,

		t.amount as original_amount,
		COALESCE(t.transfer_status, 'pending') as status,
		COALESCE(t.notes, '') as notes,
		t.created_at

	FROM transfer t
	JOIN wallets w_sender ON t.sender_wallet_id = w_sender.id
	JOIN wallets w_receiver ON t.receiver_wallet_id = w_receiver.id
	LEFT JOIN profile p_sender ON w_sender.user_id = p_sender.users_id
	LEFT JOIN profile p_receiver ON w_receiver.user_id = p_receiver.users_id
	WHERE w_sender.user_id = $1 OR w_receiver.user_id = $1
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
			&history.Type,
			&history.ProfilePicture,
			&history.ContactName,
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

	// Return error if no transactions found (optional)
	if len(histories) == 0 {
		return nil, errors.New("no transactions found")
	}

	return histories, nil
}
