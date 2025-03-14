package stores

import (
	"Stant/ECommerce/internal/domain"
	"database/sql"
)

type SellerStore struct {
	db *sql.DB
}

func NewSellerStore(db *sql.DB) *SellerStore {
	return &SellerStore{db: db}
}

const getSeller = `
  SELECT * 
  FROM sellers 
  WHERE seller_id = $1 LIMIT 1
  ;
  `

func (st SellerStore) Read(categoryID int) (domain.Seller, error) {
	row := st.db.QueryRow(getSeller, categoryID)
	seller, err := scanSeller(row)
	if err != nil {
		return seller, err
	}
	return seller, nil
}

const getSellers = `
  SELECT * 
  FROM sellers 
  ;
  `

func (st SellerStore) ReadAll() ([]domain.Seller, error) {
	rows, err := st.db.Query(getSellers)
	if err != nil {
		return nil, err
	}

	sellers := []domain.Seller{}
	defer rows.Close()
	for rows.Next() {
		seller, err := scanSeller(rows)
		if err != nil {
			return nil, err
		}
		sellers = append(sellers, seller)
	}
	return sellers, nil
}

func scanSeller(row sqlRow) (domain.Seller, error) {
	var id int32
	var name string
	if err := row.Scan(&id, &name); err != nil {
		return domain.Seller{}, err
	}
	return domain.NewSeller(id, name), nil
}
