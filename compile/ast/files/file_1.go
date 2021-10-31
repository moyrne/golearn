package files

//go:generate ast -file file_1.go -struct File1
type File1 struct {
	Name      string  `json:"name,omitempty"` // Name file1_name.
	Age       int     `json:"age"`            // file1_age
	List      []int64 `json:"list"`
	MultiList [][]int
	Sub1      File1Sub1 `json:"sub_1"` // file1_sub1
}

type File1Sub1 struct {
	Sub1Name string    `json:"sub_1_name"` // sub1_name
	Sub2     File1Sub2 `json:"sub_2"`
}

type File1Sub2 struct {
	Sub2Name string `json:"sub_2_name"`
}
