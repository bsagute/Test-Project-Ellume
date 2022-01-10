package model

//For Profile Edit
type SolarClientModel struct {
	FirstName  string `json:"FirstName" bson:"FirstName"`
	LastName   string `json:"LastName" bson:"LastName"`
	MiddleName string `json:"MiddleName" bson:"MiddleName"`
	Age        int    `json:"Age" bson:"Age"`
}

//For Profile Edit
type EllumeEmployeeModel struct {
	FirstName  string `json:"FirstName" bson:"FirstName"`
	LastName   string `json:"LastName" bson:"LastName"`
	FullName   string `json:"FullName" bson:"FullName"`
	Age        int    `json:"Age" bson:"Age"`
	BloodGroup string `json:"BloodGroup" bson:"BloodGroup"`
	Remark     string `json:"Remark" bson:"Remark"`
	MobileNo   string `json:"MobileNo" bson:"MobileNo"`
	EmailId    string `json:"EmailId" bson:"EmailId"`
}

type Branch struct {
	BranchId   string
	BranchName string
	IsEnabled  bool `json:"isEnabled" bson:"isEnabled"`
}
