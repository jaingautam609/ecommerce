package dbHelper

import (
	"database/sql"
	"ecommerce/models"
	"errors"
	"github.com/jmoiron/sqlx"
)

func AddItemType(db *sqlx.DB, Id int, types models.ItemType) error {
	SQL := `insert into item_type(item_type,added_by) values ($1,$2)`
	_, err := db.Exec(SQL, types.Type, Id)
	if err != nil {
		return err
	}
	return nil
}

func AddItem(db *sqlx.DB, item models.Item, adminId int) error {
	SQL := `insert into item(Type_id,item_name,added_by,price) values ($1,$2,$3,$4)`
	_, err := db.Exec(SQL, item.TypeId, item.Name, adminId, item.Price)
	if err != nil {
		return err
	}
	return nil
}
func DeleteItem(db *sqlx.DB, itemId int) error {
	SQL := `UPDATE item
			SET archive_at = current_timestamp
			WHERE id=$1;`
	_, err := db.Exec(SQL, itemId)
	if err != nil {
		return err
	}
	return nil
}
func CountUsers(db *sqlx.DB) (int, error) {
	var count int
	SQL := `select count(u.user_name)
				from users u join user_role
				    t on u.id=t.user_id
						where u.archive_at is null`
	err := db.QueryRowx(SQL).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}
func AllUsers(db *sqlx.DB, limit int, offset int) ([]models.Register, error) {
	var userInfo []models.Register
	SQL := `select u.user_name,u.user_email,u.joined_at ,t.user_type
				from users u join user_role
				    t on u.id=t.user_id
						where u.archive_at is null limit $1 offset $2`
	rows, err := db.Query(SQL, limit, offset)
	if err != nil {
		return userInfo, err
	}
	for rows.Next() {
		var info models.Register
		err = rows.Scan(&info.Name, &info.Email, &info.JoinedOn, &info.Type)
		if err != nil {
			return userInfo, err
		}
		userInfo = append(userInfo, info)
	}
	return userInfo, nil
}
func CountProduct(db *sqlx.DB) (int, error) {
	var count int
	SQL := `select count(item_name)
				from item 
				where archive_at is null`
	err := db.QueryRowx(SQL).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}
func AllProducts(db *sqlx.DB, limit, offset int) ([]models.Item, error) {
	var allProduct []models.Item
	SQL := `select i.item_name,i.added_by,i.added_on,i.price,array_agg(u.url)
    AS photos  from item i
                        left join item_image img
                                  on i.id=img.item_id
                        left join uploads u
                                   on u.id=img.upload_id
where i.archive_at is null GROUP BY
                               i.item_name,
                               i.added_by,
                               i.added_on,
                               i.price limit $1 offset $2;`
	err := db.Select(&allProduct, SQL, limit, offset)
	if err != nil {
		return allProduct, err
	}
	return allProduct, nil
}

func ProductById(db *sqlx.DB, id int) (models.Item, error) {
	var allProduct models.Item
	SQL := `select i.item_name,i.added_by,i.added_on,i.price,i.type_id,array_agg(u.url) 
    		AS photos 
				from item_type t
				    inner join item i 
				        on t.id = i.type_id
				    left join  item_image img 
				        on i.id=img.item_id 	
    				left join uploads u
    				    on u.id=img.upload_id  
				where i.id =$1 and i.archive_at is null GROUP BY
  				i.item_name,
  				i.added_by,
  				i.added_on,
  				i.type_id,
  				i.price;`
	err := db.Get(&allProduct, SQL, id)
	if err != nil {
		return allProduct, err
	}

	return allProduct, nil
}
func ProductByType(db *sqlx.DB, id int, limit, offset int) ([]models.Item, error) {
	var allProduct []models.Item
	SQL := `select i.item_name,i.added_by,i.added_on,i.price,array_agg(u.url) 
    		AS photos
				from item_type t
					inner join item i 
					    on t.id=i.type_id
					left join item_image img 
					    on i.id=img.item_id 
    				left join uploads u
    				    on u.id=img.upload_id
			where  t.id=$1 and i.archive_at is null GROUP BY
  				i.item_name,
  				i.added_by,
  				i.added_on,
  				i.price limit $2 offset $3;`
	err := db.Select(&allProduct, SQL, id, limit, offset)
	if err != nil {
		return allProduct, err
	}
	return allProduct, nil
}
func CountProductByType(db *sqlx.DB, id int) (int, error) {
	var count int
	SQL := `select count(i.item_name)
				from item_type t
					inner join item i 
					    on t.id=i.type_id
					left join item_image img 
					    on i.id=img.item_id 
    				left join uploads u
    				    on u.id=img.upload_id
			where  t.id=$1 and i.archive_at is null GROUP BY
  				i.item_name`
	err := db.QueryRowx(SQL).Scan(&count, id)
	if err != nil {
		return count, err
	}
	return count, nil
}
func AddToCart(db *sqlx.DB, getItem models.Item, quantity int, cartId int, itemType string, id int) error {
	SQL := `insert into cart_item(cart_id,item_name,item_type,quantity,item_id,price) values($1,$2,$3,$4,$5,$6)`
	_, err := db.Exec(SQL, cartId, getItem.Name, itemType, quantity, id, getItem.Price)
	if err != nil {
		return err
	}
	return nil
}
func IncreaseInCart(db *sqlx.DB, cartId int, id int, quantity int) error {
	SQL := `UPDATE cart_item 
				SET quantity = quantity + $3 
					WHERE item_id = $1 AND cart_id = $2`
	_, err := db.Exec(SQL, id, cartId, quantity)
	if err != nil {
		return err
	}
	return nil
}
func DeleteFromCart(db *sqlx.DB, itemId, cartId int) error {
	SQL := `DELETE from cart_item 
       			where item_id = $1
       			  AND cart_id = $2;`
	_, err := db.Exec(SQL, itemId, cartId)
	if err != nil {
		return err
	}
	return nil
}
func DecreaseFromCart(db *sqlx.DB, itemId, cartId int) error {
	SQL := `UPDATE cart_item 
				SET quantity = quantity - 1 
					WHERE item_id = $1 AND cart_id = $2`
	_, err := db.Exec(SQL, itemId, cartId)
	if err != nil {
		return err
	}
	return nil
}
func GetQuantity(db *sqlx.DB, itemId, cartId int) (int, error) {
	var quantity int
	SQL := `select quantity
				from cart_item 
					where item_id=$1 and
						cart_id=$2`
	err := db.Get(&quantity, SQL, itemId, cartId)
	if err != nil {
		if err == sql.ErrNoRows {
			return quantity, errors.New("user not found")
		}
		return quantity, err
	}
	return quantity, nil
}
func AssignCart(tx *sqlx.Tx, userId int) error {
	SQL := `insert into cart(assign_to) values($1)`
	_, err := tx.Exec(SQL, userId)
	if err != nil {
		return err
	}
	return nil
}
func GetCartId(db *sqlx.DB, userId int) (int, error) {
	var cartId int
	SQL := `SELECT id FROM cart WHERE assign_to = $1`
	err := db.Get(&cartId, SQL, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return cartId, errors.New("user not found")
		}
		return cartId, err
	}
	return cartId, nil
}
func GetType(db *sqlx.DB, typeId int) (string, error) {
	var itemType string
	SQL := `select item_type from item_type where id=$1`
	err := db.Get(&itemType, SQL, typeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return itemType, errors.New("item type not exist")
		}
		return itemType, err
	}
	return itemType, nil
}

func ShowCart(db *sqlx.DB, cartId int) ([]models.CartItem, error) {
	var allItems []models.CartItem
	SQL := `select cart_id,
       				item_name,
       				item_type,
       				quantity,
       				item_id,
       				price 
			from cart_item where cart_id=$1`
	rows, err := db.Query(SQL, cartId)
	if err != nil {
		return allItems, err
	}
	for rows.Next() {
		var allItem models.CartItem
		err = rows.Scan(&allItem.CartId, &allItem.ItemName, &allItem.ItemType, &allItem.Quantity, &allItem.ItemId, &allItem.Price)
		if err != nil {

			return allItems, err
		}
		allItems = append(allItems, allItem)
	}
	return allItems, nil
}

func Checkout(db *sqlx.DB, cartId int) error {
	SQL := `DELETE from cart_item where cart_id = $1;`
	_, err := db.Exec(SQL, &cartId)
	if err != nil {
		return err
	}
	return nil
}
func Upload(tx *sqlx.Tx, upload models.Uploads) (int, error) {
	var uploadId int
	SQL := `insert into uploads(path,name,url) values($1,$2,$3) RETURNING id`
	err := tx.QueryRowx(SQL, &upload.Path, &upload.Name, &upload.Url).Scan(&uploadId)
	if err != nil {
		return uploadId, err
	}
	return uploadId, nil
}
func ItemImage(tx *sqlx.Tx, itemId, uploadId int) error {
	SQL := `insert into item_image(item_id,upload_id) values($1,$2)`
	_, err := tx.Exec(SQL, itemId, uploadId)
	if err != nil {
		return err
	}
	return nil
}
func EnterEmail(db *sqlx.DB, email string) error {
	SQL := `insert into users(user_email,is_verified_by_email) values ($1,true)`
	_, err := db.Exec(SQL, email)
	if err != nil {
		return err
	}
	return nil
}
func EnterNumber(db *sqlx.DB, number string) error {
	SQL := `insert into users(phone_no,is_verified_by_phone) values ($1,true)`
	_, err := db.Exec(SQL, number)
	if err != nil {
		return err
	}
	return nil
}
