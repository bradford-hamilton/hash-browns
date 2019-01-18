package hashbrown

import "testing"

func TestCreate(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"hash test 1",
			args{"angryMonkey"},
			"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==",
		},
		{
			"hash test 2",
			args{"I'm another password"},
			"mdtwf8cuLCjI7QAC23+g+ohUMNpaP0mckVfYlgggAwMKxzs0bfmzD1d1k5fvNwjIbVyWLKDnfcuZOodVcnaeUQ==",
		},
		{
			"hash test 3",
			args{"And another just for fun"},
			"VyKtgs4Cg2Ab6F8zjQ9RbpNhETs4gKcfOd59A5AAW3CpadBOajAULsCVdj7USODDTU3oegOL65+bXQ8AHPlk7g==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Create(tt.args.password); got != tt.want {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
