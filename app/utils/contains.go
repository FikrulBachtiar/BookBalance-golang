package utils

func StationContains(stationCodes []string, stationCode string) bool {
	isValid := false

	for _, code := range stationCodes {
		if code == stationCode {
			isValid = true
			break
		}
	}

	return isValid
}