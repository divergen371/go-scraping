package main

import (
	"fmt"

	"gorm.io/gorm"
)

func createLatestItems(items []Item, db *gorm.DB) error {
	stmt := &gorm.Statement{DB: db}
	err := stmt.Parse(&LatestItem{})

	if err != nil {
		return fmt.Errorf("get latest_items table name error: %w", err)
	}

	if err := db.Exec("TRUNCATE " + stmt.Schema.Table).Error; err != nil {
		return fmt.Errorf("trancate latest_items table error: %w", err)
	}

	var insertRecords []LatestItem
	for _, item := range items {
		insertRecords = append(insertRecords, LatestItem{Item: item})
	}

	if err := db.CreateInBatches(insertRecords, 100).Error; err != nil {
		return fmt.Errorf("bulk insert to latest_items error: %w", err)
	}

	return nil
}

func updateItemMaster(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// INSERT
		var newItems []LatestItem
		err := tx.Unscoped().Joins("LEFT JOIN item_master ON latest_items." +
			"url = item_master.URL").Where("item_master.name IS NULL").Find(
			&newItems).Error
		if err != nil {
			return fmt.Errorf("extract for bulk insert to item_master error: %w", err)
		}

		var insertRecords []ItemMaster
		for _, newItem := range newItems {
			insertRecords = append(insertRecords, ItemMaster{Item: newItem.Item})
			fmt.Printf("Index item is created: %s\n", newItem.URL)
		}
		if err := tx.CreateInBatches(insertRecords, 100).Error; err != nil {
			return fmt.Errorf("bulk insert to item_master error: %w", err)
		}

		// UPDATE
		var updateItems []LatestItem
		err = tx.Unscoped().Joins("INNER JOIN item_master ON latest_items." +
			"url = item_master.url").Where("latest_items.name <> item_master." +
			"name OR latest_items.price <> item_master." +
			"price OR item_master.deleted_at IS NOT NULL").Find(&updateItems).Error
		if err != nil {
			return fmt.Errorf("update error: %w", err)
		}
		for _, updateItem := range updateItems {
			err := tx.Unscoped().Model(ItemMaster{}).Where("url = ?",
				updateItem.URL).Updates(
				map[string]interface{}{
					"name":  updateItem.Name,
					"price": updateItem.Price, "deleted_at": nil,
				}).Error
			if err != nil {
				return fmt.Errorf("update error: %w", err)
			}
			fmt.Printf("Index item is updated: %s\n", updateItem.URL)
		}

		// DELETE
		var deletedItems []ItemMaster
		if err := tx.Where("NOT EXISTS(SELECT 1 FROM latest_items li WHERE li.url = item_master.url)").
			Find(&deletedItems).Error; err != nil {
			return fmt.Errorf("delete error: %w", err)
		}
		var ids []uint
		for _, deleteItem := range deletedItems {
			ids = append(ids, deleteItem.ID)
			fmt.Printf("Index item is deleted: %s\n", deleteItem.URL)
		}
		if len(ids) > 0 {
			if err := tx.Delete(&deletedItems).Error; err != nil {
				return fmt.Errorf("delete error: %w", err)
			}
		}
		return nil
	})
}
