package views

const (
	// AlertLvlError turns the alert red
	AlertLvlError = "danger"
	// AlertLvlWarning turns the alert yellow
	AlertLvlWarning = "warning"
	// AlertLvlInfo turns the alert blue
	AlertLvlInfo = "info"
	// AlertLvlSuccess turns the alert green
	AlertLvlSuccess = "success"
	//AlertMsgGeneric is just a generic error msg we show
	AlertMsgGeneric = "Something went wrong. Please try again and do contact us if the problem persists."
)

//Alert is the data passed in to the Alert
//it has two fields as the level one denotes what type of msg we sho
type Alert struct {
	Level   string
	Message string
}

// Data is the generic placeholder. It gets the Alert struct +
// it provides the Yield data for the generic template
type Data struct {
	Alert *Alert
	Yield interface{}
}
