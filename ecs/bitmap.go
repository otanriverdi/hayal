package ecs

import "math/bits"

func buildBitmap(cmpIds ...componentId) bitmap {
	var bitmap bitmap
	for _, cmpId := range cmpIds {
		bitmap[cmpId/64] |= 1 << (cmpId % 64)
	}
	return bitmap
}

func setBitmap(bitmap bitmap, cmpId componentId) bitmap {
	bitmap[cmpId/64] |= 1 << (cmpId % 64)
	return bitmap
}

func clearBitmap(bitmap bitmap, cmpId componentId) bitmap {
	bitmap[cmpId/64] &^= 1 << (cmpId % 64)
	return bitmap
}

func bitmapIsSubset(a, b bitmap) bool {
	for key, aValue := range a {
		bValue := b[key]
		if aValue&^bValue != 0 {
			return false
		}
	}
	return true
}

func extractBitmapCmps(b bitmap) []componentId {
	var cmpIds []componentId
	for index, word := range b {
		tempWord := word
		for tempWord != 0 {
			tz := bits.TrailingZeros64(tempWord)
			cmpId := componentId(uint32(index)*64 + uint32(tz))
			cmpIds = append(cmpIds, cmpId)
			tempWord &= tempWord - 1
		}
	}
	return cmpIds
}
