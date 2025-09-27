package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Belalai-E-Wallet-Backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChartRepository struct {
	db *pgxpool.Pool
}

func NewChartRepository(db *pgxpool.Pool) *ChartRepository {
	return &ChartRepository{db: db}
}

func (cr *ChartRepository) GetChartData(c context.Context, user_id int, filter string) (models.ChartData, error) {

	sqlWeekData := `WITH date_series AS (
			SELECT generate_series(
					CURRENT_DATE - INTERVAL '6 days',
					CURRENT_DATE,
					'1 day'
			)::DATE AS date
	),
	transactions AS (
			SELECT id AS wallet_id FROM wallets WHERE user_id = $1
	),
	daily_data AS (
			SELECT
					t.created_at::DATE AS date, t.amount AS income, 0 AS expense
			FROM topup t JOIN wallets_topup wt ON t.id = wt.topup_id
			WHERE wt.wallets_id = (SELECT wallet_id FROM transactions) AND t.topup_status = 'success' AND t.created_at >= CURRENT_DATE - INTERVAL '6 days'
			UNION ALL
			SELECT
					tr.created_at::DATE AS date, tr.amount AS income, 0 AS expense
			FROM transfer tr
			WHERE tr.receiver_wallet_id = (SELECT wallet_id FROM transactions) AND tr.transfer_status = 'success' AND tr.created_at >= CURRENT_DATE - INTERVAL '6 days'
			UNION ALL
			SELECT
					tr.created_at::DATE AS date, 0 AS income, tr.amount AS expense
			FROM transfer tr
			WHERE tr.sender_wallet_id = (SELECT wallet_id FROM transactions) AND tr.transfer_status = 'success' AND tr.created_at >= CURRENT_DATE - INTERVAL '6 days'
	),
	aggregated_daily AS (
			SELECT
					ds.date AS label,
					COALESCE(SUM(dd.income), 0) AS income,
					COALESCE(SUM(dd.expense), 0) AS expense
			FROM date_series ds
			LEFT JOIN daily_data dd ON ds.date = dd.date
			GROUP BY ds.date
			ORDER BY ds.date
	)
	SELECT
			array_agg(TO_CHAR(label, 'YYYY-MM-DD')) AS labels,
			array_agg(income) AS income_data,
			array_agg(expense) AS expense_data
	FROM
			aggregated_daily`

	sqlMonthData := `WITH week_series AS (
			SELECT DISTINCT DATE_TRUNC('week', generate_series(
					CURRENT_DATE - INTERVAL '29 days',
					CURRENT_DATE,
					'1 day'
			))::DATE AS week_start
	),
	transactions AS (
			SELECT id AS wallet_id FROM wallets WHERE user_id = $1
	),
	weekly_data AS (
			SELECT
					DATE_TRUNC('week', t.created_at)::DATE AS week_start,
					t.amount AS income,
					0 AS expense
			FROM topup t
			JOIN wallets_topup wt ON t.id = wt.topup_id
			WHERE wt.wallets_id = (SELECT wallet_id FROM transactions) AND t.topup_status = 'success' 
			AND t.created_at >= CURRENT_DATE - INTERVAL '29 days' -- Batas 30 hari

			UNION ALL

			SELECT
					DATE_TRUNC('week', tr.created_at)::DATE AS week_start,
					tr.amount AS income,
					0 AS expense
			FROM transfer tr
			WHERE tr.receiver_wallet_id = (SELECT wallet_id FROM transactions) AND tr.transfer_status = 'success' 
			AND tr.created_at >= CURRENT_DATE - INTERVAL '29 days' -- Batas 30 hari

			UNION ALL

			SELECT
					DATE_TRUNC('week', tr.created_at)::DATE AS week_start,
					0 AS income,
					tr.amount AS expense
			FROM transfer tr
			WHERE tr.sender_wallet_id = (SELECT wallet_id FROM transactions) AND tr.transfer_status = 'success' 
			AND tr.created_at >= CURRENT_DATE - INTERVAL '29 days' -- Batas 30 hari
	),
	aggregated_weekly AS (
			SELECT
					ws.week_start AS label_date,
					COALESCE(SUM(wd.income), 0) AS income,
					COALESCE(SUM(wd.expense), 0) AS expense
			FROM week_series ws
			LEFT JOIN weekly_data wd ON ws.week_start = wd.week_start
			GROUP BY ws.week_start
			ORDER BY ws.week_start
	)
	SELECT
			array_agg(TO_CHAR(label_date, 'YYYY-MM-DD')) AS labels,
			array_agg(income) AS income_data,
			array_agg(expense) AS expense_data
	FROM
			aggregated_weekly;`

	sqlYearData := `WITH month_series AS (
			SELECT generate_series(
					DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '11 months',
					DATE_TRUNC('month', CURRENT_DATE),
					'1 month'
			)::DATE AS month_start
	),
	transactions AS (
			SELECT id AS wallet_id FROM wallets WHERE user_id = $1
	),
	monthly_data AS (
			SELECT
					DATE_TRUNC('month', t.created_at)::DATE AS month_start, t.amount AS income, 0 AS expense
			FROM topup t JOIN wallets_topup wt ON t.id = wt.topup_id
			WHERE wt.wallets_id = (SELECT wallet_id FROM transactions) AND t.topup_status = 'success' AND t.created_at >= DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '11 months'
			UNION ALL
			SELECT
					DATE_TRUNC('month', tr.created_at)::DATE AS month_start, tr.amount AS income, 0 AS expense
			FROM transfer tr
			WHERE tr.receiver_wallet_id = (SELECT wallet_id FROM transactions) AND tr.transfer_status = 'success' AND tr.created_at >= DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '11 months'
			UNION ALL
			SELECT
					DATE_TRUNC('month', tr.created_at)::DATE AS month_start, 0 AS income, tr.amount AS expense
			FROM transfer tr
			WHERE tr.sender_wallet_id = (SELECT wallet_id FROM transactions) AND tr.transfer_status = 'success' AND tr.created_at >= DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '11 months'
	),
	aggregated_monthly AS (
			SELECT
					ms.month_start AS label,
					COALESCE(SUM(md.income), 0) AS income,
					COALESCE(SUM(md.expense), 0) AS expense
			FROM month_series ms
			LEFT JOIN monthly_data md ON ms.month_start = md.month_start
			GROUP BY ms.month_start
			ORDER BY ms.month_start
	)
	SELECT
			array_agg(TO_CHAR(label, 'YYYY-MM')) AS labels,
			array_agg(income) AS income_data,
			array_agg(expense) AS expense_data
	FROM
			aggregated_monthly;`

	var queryExec string
	switch filter {
	case "seven_days":
		queryExec = sqlWeekData
	case "five_weeks":
		queryExec = sqlMonthData
	case "twelve_months":
		queryExec = sqlYearData
	default:
	}

	var chartData models.ChartData
	if err := cr.db.QueryRow(c, queryExec, user_id).Scan(&chartData.Labels, &chartData.IncomeData, &chartData.ExpenseData); err != nil {
		if err == sql.ErrNoRows {
			return models.ChartData{}, fmt.Errorf("no data found for user %d", user_id)
		}
		return models.ChartData{}, err
	}
	// if no error return data, and error is nil
	return chartData, nil
}
