package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type WarehouseId int

func (w *WarehouseId) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		id, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		*w = WarehouseId(id)
		return nil
	}

	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*w = WarehouseId(num)
		return nil
	}

	return fmt.Errorf("не удалось распарсить warehouse_id")
}

func (w WarehouseId) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(w)))
}
