package character

type weapons struct {
	Name         string
	Encumberance int
	Group        string
	Damage       string
	MinRange     int
	MaxRange     int
	Reload       string
}

type armour struct {
	Name             string
	ArmourType       string
	Encumberance     int
	LocationsCovered string
	ArmourPoints     int
}
