package main

type mainDats struct {
	TitleContext   string
	LinkToContext  string
	Contenders     []contenders
	Jurys          []jury
	RewardsSumCont int64
}
type contenders struct {
	IDS          int64
	Address      string
	AverageScore float64
	GovermentD   *goverment
	Reject       int64
	Jury         bool
}

type jury struct {
	Address    string
	PublicKey  string
	VFor       int
	VAbstained int
	VAgainst   int
}

type goverment struct {
	JurorsAbstained []string `json:"jurorsAbstained"`
	JurorsAgainst   []string `json:"jurorsAgainst"`
	JurorsFor       []string `json:"jurorsFor"`
	Marks           []string `json:"marks"`
}

type votes struct {
	JuryFor       int64
	JuryAbstained int64
	JuryAgainst   int64
}

type resultContenders struct {
	Addresses []string `json:"addresses"`
	Ids       []string `json:"ids"`
}

type resContestInfo struct {
	JuryAddresses []string `json:"juryAddresses"`
	JuryKeys      []string `json:"juryKeys"`
	Link          string   `json:"link"`
	Title         string   `json:"title"`
	Hash          string   `json:"hash"`
}

type req struct {
	ID int64 `json:"id"`
}

type marksAddress struct {
	Address string
	Mark    int64
}

type allVotes struct {
	Mark    int64
	Abstain bool
}

type statJurys struct {
	Efficiency      float64
	CountEfficiency float64
	AvgDiff         float64
	CountAvgDiff    float64
}

type dataForSort []statJurys

type dataSlice []resultTable

type resultTable struct {
	address    string
	efficiency float64
	avgDiff    float64
}

const (
	linkToExplorer = "https://ton.live/accounts?section=details&id="
)

var resultData map[string]statJurys

func init() {
	resultData = make(map[string]statJurys)
}

// Len is part of sort.Interface.
func (d dataSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d dataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d dataSlice) Less(i, j int) bool {
	return d[i].efficiency > d[j].efficiency
}
