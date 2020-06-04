package dynamic

var (
	HomeBody = `
<table style="width:100%">
  <tr>
	<th>ID</th>
	<th>Name</th>
	<th>Progress</th>
	<th>Status</th>
  </tr>
{{range .}}
  <tr>
	<td>{{.Id}}</td>
	<td>{{.Name}}</td>
	<td>{{.Percent}}%</td>
	<td>{{.Status}}</td>
  </tr>
{{end}}
</table>
`
)


type Entry struct {
	Username string
	Password string
	Secret   string
	Projects []Project
}

type Project struct {
	Id		uint64
	Name	string
	Percent uint64
	Status	string
}

type StatCode int

const (
	SUCCEEDED = iota; FAILED; PENDING; BUILDING;
)

func (s StatCode) String() string {
	return [...]string { "SUCCEEDED", "FAILED", "PENDING", "BUILDING" }[s]
}
