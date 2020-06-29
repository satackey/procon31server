package apispec

// Matches は 事前情報 の集まりです
type Matches []*Match

// Match は 事前情報 を表します
type Match struct {
	ID             int    `json:"id"`
	IntervalMillis int    `json:"intervalMillis"`
	MatchTo        string `json:"matchTo"`
	TeamID         int    `json:"teamID"`
	TurnMillis     int    `json:"turnMillis"`
	Turns          int    `json:"turns"`
}

// FieldStatus は 試合状態を表します
type FieldStatus struct {
	Width             int                 `json:"width"`
	Height            int                 `json:"height"`
	Points            [][]int             `json:"points"`
	StartedAtUnixtime int64               `json:"startedAtUnixTime"`
	Turn              int                 `json:"turn"`
	Cells             [][]Cell            `json:"cell"`
	Teams             []Team              `json:"teams"`
	Actions           []FieldStatusAction `json:"actions"`
}

// Cell は セルの情報 を表します
type Cell struct {
	Status string
	TeamID int
}

// Team は team情報 を表します
type Team struct {
	TeamID    int     `json:"teamID"`
	Agents    []Agent `json:"agents"`
	WallPoint int     `json:"wallPoint"`
	AreaPoint int     `json:"areaPoint"`
}

// Agent は エージェント情報 を表します
type Agent struct {
	AgentID int `json:"agentID"`
	X       int `json:"x"`
	Y       int `json:"y"`
}

// FieldStatusAction は エージェントの移動情報 を表します
type FieldStatusAction struct {
	AgentID int    `json:"agentID"`
	DX      int    `json:"dx"`
	DY      int    `json:"dy"`
	Type    string `json:"type"`
	Apply   int    `json:"apply"`
	Turn    int    `json:"turn"`
}

// Update は 行動更新情報 の集まりです
type Update struct {
	Actions []UpdateAction `json:"actions"`
}

// UpdateAction は行動情報 を表します
type UpdateAction struct {
	AgentID int    `json:"agentID"`
	DX      int    `json:"dx"`
	DY      int    `json:"dy"`
	Type    string `json:"type"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}

// Ping は 動作確認 をします
type Ping struct {
	Status string `json:"status"`
}

// Match401 は
type Match401 struct {
	Status string `json:"status"`
}

// Field400 は
type Field400 struct {
	StartAtUnixTime int64  `json:"startAtUnixTime"`
	Status          string `json:"status"`
}

// Field401 は
type Field401 struct {
	Status string `json:"status"`
}

// Update400 は
type Update400 struct {
	StartAtUnixTime int64  `json:"startAtUnixTime"`
	Status          string `json:"status"`
}

//Update401 は
type Update401 struct {
	Status string `json:"status"`
}

//Ping401 は
type Ping401 struct {
	Status string `json:"status"`
}
