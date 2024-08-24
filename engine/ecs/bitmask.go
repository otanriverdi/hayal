package ecs

func setBit(mask *TypeId, bit ComponentId) {
	*mask |= (1 << bit)
}

func clearBit(mask *TypeId, bit ComponentId) {
	*mask &= ^(1 << bit)
}

func cmpBit(mask TypeId, bit ComponentId) bool {
	return (mask & (1 << bit)) != 0
}
