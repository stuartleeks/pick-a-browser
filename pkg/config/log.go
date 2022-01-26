package config

type Log struct {
	Level string
}

func parseLog(rootNode map[string]interface{}) (Log, error) {
	logNode, err := getObjectChildNode(rootNode, "log", false)
	if err != nil {
		return Log{}, err
	}

	level, err := parseLogLevel(logNode)
	if err != nil {
		return Log{}, err
	}

	return Log{
		Level: level,
	}, nil
}

func parseLogLevel(parentNode map[string]interface{}) (string, error) {
	logLevel, err := getOptionalString(parentNode, "level")
	if err != nil {
		return "", err
	}

	if logLevel == nil {
		return "", nil
	}
	return *logLevel, nil
}
