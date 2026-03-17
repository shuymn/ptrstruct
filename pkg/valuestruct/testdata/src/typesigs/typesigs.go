package typesigs

type User struct {
	Name string
}

// Interface method parameter by pointer: NG.
type Saver interface {
	Save(u *User) // want `interface method Save parameter u uses pointer to struct User; use User`
}

// Function type result by pointer: NG.
type LoadFunc func() *User // want `function type LoadFunc result uses pointer to struct User; use User`

// Named container type with pointer elements: NG.
type Users []*User // want `named type Users uses slice element User by pointer`
