package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/GoncharovFyodor/hezzltest/internal/domain"
	"github.com/GoncharovFyodor/hezzltest/internal/models"
	_ "github.com/lib/pq"
)

type GoodsPostgres struct {
	db *sql.DB
}

func NewGoodsRepository(db *sql.DB) *GoodsPostgres {
	return &GoodsPostgres{db: db}
}

func (r *GoodsPostgres) Get(ctx context.Context, limit, offset int) (domain.GoodList, error) {
	// Запрос для получения товаров
	q := `SELECT id, project_id, name, description, priority, removed, created_at FROM goods LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return domain.GoodList{}, fmt.Errorf("не удалось отправить запрос для получения товаров: %w", err)
	}
	defer rows.Close()

	var goods []models.Good
	for rows.Next() {
		var good models.Good
		err := rows.Scan(
			&good.ID,
			&good.ProjectID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreatedAt,
		)
		if err != nil {
			return domain.GoodList{}, fmt.Errorf("не удалось прочитать список товаров: %w", err)
		}
		goods = append(goods, good)
	}

	if err := rows.Err(); err != nil {
		return domain.GoodList{}, fmt.Errorf("ошибка получения строк: %w", err)
	}

	// Запрос для получения количества удаленных товаров
	var removedCount int
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(id) FROM goods WHERE removed = true`).Scan(&removedCount)
	if err != nil {
		return domain.GoodList{}, fmt.Errorf("не удалось получить количество удаленных товаров: %w", err)
	}

	return domain.GoodList{
		Meta: domain.Meta{
			Total:   len(goods),
			Removed: removedCount,
			Limit:   limit,
			Offset:  offset,
		},
		Goods: goods,
	}, nil
}

func (r *GoodsPostgres) Create(ctx context.Context, projectID int, input domain.CreateGoodRequest) (models.Good, error) {
	var created models.Good
	q := `INSERT INTO goods (project_id, name) 
VALUES ($1, $2)
RETURNING id, project_id, name, description, priority, removed, created_at`

	err := r.db.QueryRowContext(ctx, q, projectID, input.Name).Scan(
		&created.ID,
		&created.ProjectID,
		&created.Name,
		&created.Description,
		&created.Priority,
		&created.Removed,
		&created.CreatedAt,
	)
	if err != nil {
		return models.Good{}, fmt.Errorf("не удалось создать товар: %w", err)
	}
	return created, nil
}

func (repo *GoodsPostgres) Update(ctx context.Context, projectID, id int, input domain.UpdateGoodRequest) (models.Good, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Good{}, fmt.Errorf("не удалось запустить транзакцию: %w", err)
	}
	defer tx.Rollback()

	// Блокировка записи для чтения
	var current models.Good
	err = tx.QueryRowContext(ctx, `
		SELECT id, project_id, name, description, priority, removed, created_at 
		FROM goods 
		WHERE id = $1 AND project_id = $2
		FOR UPDATE`,
		id, projectID).Scan(
		&current.ID,
		&current.ProjectID,
		&current.Name,
		&current.Description,
		&current.Priority,
		&current.Removed,
		&current.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Good{}, fmt.Errorf("товар не найден")
		}
		return models.Good{}, fmt.Errorf("не удалось поставить блокировку на чтение записи: %w", err)
	}

	// Обновление
	var updated models.Good
	err = tx.QueryRowContext(ctx, `
		UPDATE goods 
		SET 
			name = $1, 
			description = COALESCE(NULLIF($2, ''), description) 
		WHERE id = $3 AND project_id = $4
		RETURNING id, project_id, name, description, priority, removed, created_at`,
		input.Name, input.Description, id, projectID).Scan(
		&updated.ID,
		&updated.ProjectID,
		&updated.Name,
		&updated.Description,
		&updated.Priority,
		&updated.Removed,
		&updated.CreatedAt,
	)

	if err != nil {
		return models.Good{}, fmt.Errorf("не удалось обновить товар: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return models.Good{}, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	return updated, nil
}

func (r *GoodsPostgres) Delete(ctx context.Context, projectID, ID int) (models.DeletedGood, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.DeletedGood{}, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	// Блокировка записи для чтения
	var current models.DeletedGood
	query := `SELECT id, project_id, removed 
	          FROM goods WHERE id = $1 AND project_id = $2 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, ID, projectID).Scan(
		&current.ID,
		&current.ProjectID,
		&current.Removed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DeletedGood{}, err
		}
		return models.DeletedGood{}, fmt.Errorf("не удалось получить товар и совершить блокировку:: %w", err)
	}

	// Пометка товара как удаленного
	_, err = tx.ExecContext(ctx, "UPDATE goods SET removed = true WHERE id = $1 AND project_id = $2 RETURNING id, project_id, removed", ID, projectID)
	if err != nil {
		return models.DeletedGood{}, fmt.Errorf("не удалось пометить товар как удаленный: %w", err)
	}

	current.Removed = true

	if err = tx.Commit(); err != nil {
		return models.DeletedGood{}, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	return current, nil
}

func (r *GoodsPostgres) Reprioritize(ctx context.Context, projectID, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable, // Уровень изоляции для предотвращения аномалий
	})
	if err != nil {
		return models.GoodPriorities{}, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var priority int
	err = tx.QueryRowContext(ctx,
		`SELECT priority FROM goods WHERE id = $1 AND project_id = $2 FOR UPDATE`, // Блокировка FOR UPDATE
		id, projectID).Scan(&priority)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.GoodPriorities{}, fmt.Errorf("товар не найден")
		}
		return models.GoodPriorities{}, fmt.Errorf("ошибка получения приоритета: %w", err)
	}

	// Блокировка всех записей для обновления
	rows, err := tx.QueryContext(ctx,
		`SELECT id, priority FROM goods 
         WHERE priority >= $1 
         ORDER BY id ASC
         FOR UPDATE SKIP LOCKED`, // Блокировка с пропуском уже заблокированных
		priority)
	if err != nil {
		return models.GoodPriorities{}, fmt.Errorf("ошибка запроса товаров: %w", err)
	}
	defer rows.Close()

	var goods []models.Priorities
	for rows.Next() {
		var good models.Priorities
		if err := rows.Scan(&good.ID, &good.Priority); err != nil {
			return models.GoodPriorities{}, fmt.Errorf("ошибка чтения данных товара: %w", err)
		}
		goods = append(goods, good)
	}
	if err := rows.Err(); err != nil {
		return models.GoodPriorities{}, fmt.Errorf("ошибка обработки результатов: %w", err)
	}

	var updatedGoods models.GoodPriorities
	currentPriority := input.NewPriority

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE goods SET priority = $1 WHERE id = $2 RETURNING id, priority`)
	if err != nil {
		return models.GoodPriorities{}, fmt.Errorf("ошибка подготовки запроса: %w", err)
	}
	defer stmt.Close()

	for _, good := range goods {
		var updatedGood models.Priorities
		err = stmt.QueryRowContext(ctx, currentPriority, good.ID).Scan(
			&updatedGood.ID, &updatedGood.Priority)
		if err != nil {
			return models.GoodPriorities{}, fmt.Errorf("ошибка обновления товара %d: %w", good.ID, err)
		}
		updatedGoods.Priorities = append(updatedGoods.Priorities, updatedGood)
		currentPriority++
	}

	if err = tx.Commit(); err != nil {
		return models.GoodPriorities{}, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return updatedGoods, nil
}
