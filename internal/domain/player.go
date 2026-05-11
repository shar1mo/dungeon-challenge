package domain

type Player struct {
	ID int

	Registered bool
	Started    bool
	Finished   bool
	State      State

	Health int

	CurrentFloor int
	OnBossFloor  bool

	EnteredAt  int
	FinishedAt int

	FloorEnterTime       int
	CurrentFloorKills    int
	ClearedFloors        int
	TotalFloorClearTime  int
	BossEnterTime        int
	BossKillDuration     int
	BossKilled           bool
}

func NewPlayer(id int) *Player {
	return &Player{
		ID: id,
		Health: 100,
		State: StateFail,
	}
}