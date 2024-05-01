package cmd

const (
	DefaultRefactorPrompt = "Optimize the code for better performance and readability."
	RefactorPromptPattern = "Analyze the following code and return only the refactored or optimized code based on this instruction: '%s'. " +
		"Provide the refactored version only, without extra text or unchanged code.\n\n```%s```"
	BuildCommandPrompt = "Provide the simplest command line required to build the listed files. The command must be in a single line and contain no extra text or commentary:\n"
	TestCommandPrompt  = "Provide the simplest command line required to test the listed files. The command must be in a single line and contain no extra text or commentary:\n"
	LintCommandPrompt  = "Provide the simplest command line required to lint the listed files. The command must be in a single line and contain no extra text or commentary:\n"
	DocsCommandPrompt  = "Provide the simplest command line required to generate documentation for the listed files. The command must be in a single line and contain no extra text or commentary:\n"
)
