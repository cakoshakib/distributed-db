package storage

type User struct {
	tables map[string]UserTable
}

type UserTable struct {
	data map[string]interface{}
}
