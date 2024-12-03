package fixture

type NoFields struct {
}

type OneField struct {
	Field string
}

type IgnoresPrivateFields struct {
	private string
	Public  string
}
