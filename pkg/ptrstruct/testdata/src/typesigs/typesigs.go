package typesigs

type User struct {
	Name string
}

// Interface method parameter by value: NG.
type Saver interface {
	Save(u User) // want `interface method Save parameter u uses value struct User; use \*User`
}

// Function type parameter by value: NG.
type SaveFunc func(u User) // want `function type SaveFunc parameter u uses value struct User; use \*User`

// Named container type with value elements: NG.
type Users []User // want `named type Users uses slice element User by value`
