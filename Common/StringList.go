package Common

type StringList []string

func (this StringList) UniqueAdd(token string) StringList {
	for _, v := range this {
		if v == token {
			return this
		}
	}

	for i, v := range this {
		if v == "" {
			this[i] = token
			return this
		}
	}
	return append(this, token)
}

func (this StringList) Delete(token string) int {
	count := 0
	for i, v := range this {
		if v == token {
			this[i] = ""
			count++
		}
	}
	return count
}

func (this StringList) IsEmpty() bool {
	for _, v := range this {
		if v != "" {
			return false
		}
	}
	return true
}

func (this StringList) Count() int {
	count := 0
	for _, v := range this {
		if v != "" {
			count++
		}
	}
	return count
}
