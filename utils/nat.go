package utils

func GetMaxSupportUdpNat(nat int8) int8 {
	switch {
	case nat < 0:
		return -1
	case nat < 3:
		return 4
	case nat == 3:
		return 3
	case nat > 3:
		return 2
	}

	return -1
}
