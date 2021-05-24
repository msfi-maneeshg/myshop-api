package front

import (
	"database/sql"
	"myshop-api/api/data"
)

func CheckUserLoginDetails(userID, password string) (string, error) {
	var userName sql.NullString
	sqlStr := "SELECT user_name FROM `user` WHERE email_id = ? AND password = ?"
	err := data.DemoDB.QueryRow(sqlStr, userID, password).Scan(&userName)
	if err != nil && err != sql.ErrNoRows {
		return userName.String, err
	}
	return userName.String, nil
}
