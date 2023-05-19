package goutils

func QuickLoad() string {
	env := LoadEnv()
	EnableLogrus()
	LoadLocation()
	return env
}
