/*
Copyright © 2024 Chris Greaves cjgreaves97@hotmail.co.uk

See the file COPYING in the root of this repository for details.
*/
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
