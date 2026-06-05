package order

func normalizeCheckoutItems(items []checkoutLineInput) ([]checkoutLineInput, error) {
	if len(items) == 0 {
		return nil, ErrCheckoutItemEmpty
	}

	quantitiesBySku := make(map[string]int64, len(items))
	itemBySku := make(map[string]checkoutLineInput, len(items))
	orderedSkuIDs := make([]string, 0, len(items))

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrCheckoutItemInvalid
		}

		skuID := item.SkuID.String()
		if _, exists := itemBySku[skuID]; !exists {
			itemBySku[skuID] = item
			orderedSkuIDs = append(orderedSkuIDs, skuID)
		}
		quantitiesBySku[skuID] += item.Quantity
	}

	normalized := make([]checkoutLineInput, 0, len(orderedSkuIDs))
	for _, skuID := range orderedSkuIDs {
		item := itemBySku[skuID]
		item.Quantity = quantitiesBySku[skuID]
		normalized = append(normalized, item)
	}

	return normalized, nil
}
