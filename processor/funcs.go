package processor

func CreateFuncsFromModel(m Model) map[string]interface{} {
	funcs := make(map[string]interface{})

	funcs["var"] = func(key string) string {
		return m.Var[key]
	}

	funcs["env"] = func(key string) string {
		return m.Env[key]
	}

	return funcs
}
