package magiccfg

import "strings"

func formEnvName(name string) string {
	var sb strings.Builder
	runes := []rune(name)
	size := len(runes)
	for idx, rn := range runes {

		if isCapital(rn) &&
			idx > 0 &&
			idx < size-1 &&
			!isCapital(runes[idx+1]) &&
			isCapital(runes[idx-1]) {
			sb.WriteString("_")
		}

		sb.WriteRune(rn)
		if !isCapital(rn) &&
			idx < size-1 &&
			isCapital(runes[idx+1]) {
			sb.WriteString("_")
		}
	}

	return strings.ToUpper(sb.String())
}

func convertNameToEnv(name string, parentName *string, prefix string, nameMap map[string]int) string {
	envName := formEnvName(name)
	sb := strings.Builder{}
	if prefix != "" {
		sb.WriteString(strings.ToUpper(prefix))
		if !strings.HasSuffix(prefix, "_") {
			sb.WriteString("_")
		}
	}

	if num, ok := nameMap[envName]; ok && num > 1 && parentName != nil {
		sb.WriteString(*parentName)
		sb.WriteString("_")
	}

	sb.WriteString(envName)
	return sb.String()
}

func convertNameToCommandLine(name string, parentName *string) string {
	var sb strings.Builder

	if parentName != nil {
		sb.WriteString(*parentName)
		sb.WriteString(".")
	}

	runes := []rune(name)
	size := len(runes)
	for idx, rn := range runes {

		if isCapital(rn) &&
			idx > 0 &&
			idx < size-1 &&
			!isCapital(runes[idx+1]) &&
			isCapital(runes[idx-1]) {
			sb.WriteString("-")
		}

		sb.WriteRune(rn)
		if !isCapital(rn) &&
			idx < size-1 &&
			isCapital(runes[idx+1]) {
			sb.WriteString("-")
		}
	}
	return strings.ToLower(sb.String())
}

func convertNameToYAML(name string) string {
	var sb strings.Builder

	runes := []rune(name)
	size := len(runes)
	for idx, rn := range runes {

		if isCapital(rn) &&
			idx > 0 &&
			idx < size-1 &&
			!isCapital(runes[idx+1]) &&
			isCapital(runes[idx-1]) {
			sb.WriteString("-")
		}

		sb.WriteRune(rn)
		if !isCapital(rn) &&
			idx < size-1 &&
			isCapital(runes[idx+1]) {
			sb.WriteString("-")
		}
	}
	return strings.ToLower(sb.String())
}
