package nm

type Symbol struct {
	Name    string
	Version string
	Library string
}

func defaultIfEmpty(s, defaultString string) string {
	if s == "" {
		return ""
	}
	return defaultString
}
