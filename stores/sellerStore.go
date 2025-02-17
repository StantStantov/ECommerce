package stores

import (
	"Stant/ECommerce/domain"
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
