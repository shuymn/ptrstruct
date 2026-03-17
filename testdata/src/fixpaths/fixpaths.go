package fixpaths

type User struct {
	Name string
}

func MapKey(m map[User]string) {} // want `parameter m uses map key User by value`

func Array(items [2]User) {} // want `parameter items uses array element User by value`

func Stream(ch chan User) {} // want `parameter ch uses chan element User by value`

type LookupWrap struct {
	Lookup map[User]string // want `field Lookup uses map key User by value`
}

type BatchWrap struct {
	Batch [2]User // want `field Batch uses array element User by value`
}

type QueueWrap struct {
	Queue chan User // want `field Queue uses chan element User by value`
}

type InlineWrap struct {
	Inline struct{ ID string } // want `field Inline uses value struct struct{...}; use \*struct{...}`
}
