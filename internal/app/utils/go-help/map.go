package help

// M represents JSON response body.
type M map[string]interface{}

// GetInt returns int from the map for the given key.
func (m M) GetInt(key string) (int, error) {
	val, ok := m[key]
	if !ok {
		return 0, ErrMapKeyDoesNotExist
	}
	n, ok := val.(int)
	if !ok {
		return 0, ErrUnknownMapValueType
	}

	return n, nil
}

// GetString returns string from the map for the given key.
func (m M) GetString(key string) (string, error) {
	val, ok := m[key]
	if !ok {
		return "", ErrMapKeyDoesNotExist
	}
	s, ok := val.(string)
	if !ok {
		return "", ErrUnknownMapValueType
	}

	return s, nil
}
