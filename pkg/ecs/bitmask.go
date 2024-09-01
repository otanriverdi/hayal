package ecs

func setBit(mask *archetypeId, bit componentId) {
	*mask |= (1 << bit)
}

func clearBit(mask *archetypeId, bit componentId) {
	*mask &= ^(1 << bit)
}

func cmpBit(mask archetypeId, bit componentId) bool {
	return (mask & (1 << bit)) != 0
}
