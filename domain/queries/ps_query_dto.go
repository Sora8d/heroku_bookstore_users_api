package queries

type PsQuery struct {
	Equals []FieldValue `json:"equals"`
	Date   []DateValue  `json:"date"`
}

type FieldValue struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type DateValue struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
}

//This is all to obfuscate the type of db you are using and making it easier to "change" if needed, and it makes it a common library.
/* Example
{
	"equals": [
	{
		"field": "status",
		"value": "pending"
	},
	{
		"field": "seller",
		"value": 1
	}
]
}

{
	"any_equals": [
	{
		"field": "status",
		"value": "pending"
	},
	{
		"field": "seller",
		"value": 1
	}
]
}
*/
