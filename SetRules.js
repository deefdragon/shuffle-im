//This function runs through the list of rules that 
function myFunction() {
	var sheet = SpreadsheetApp.getActiveSheet();
	var rules = []
	var maxColor = "#228B22"
	var minColor = "#ffffff"
	for (var i = 1; i <= 550; i++) {
		var rangeString = `A${i}:AZ${i}`
		var range = sheet.getRange(rangeString);

		var rule = SpreadsheetApp.newConditionalFormatRule()
			.setGradientMaxpoint(maxColor)
			.setGradientMinpoint(minColor)
			.setRanges([range])
			.build();
		Logger.log("building rule at " + rangeString);
		rules.push(rule);

	}
	Logger.log("Setting Rules")
	sheet.setConditionalFormatRules(rules);
	Logger.log("Rules Set")
}