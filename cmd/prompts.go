package cmd

const (
    DefaultRefactorPrompt = "Optimize the code for better performance and readability."
    RefactorPromptPattern = "Analyze the following code and return only the refactored or optimized code based on this instruction: '%s'. " +
        "Provide the refactored version only, without extra text or unchanged code.\n\n```%s```"
    BuildCommandPrompt = "Based on the following file list, provide the simplest command line required to build these files. The command must be in a single line and contain no extra text or commentary:\n"
    TestCommandPrompt = "Based on the following file list, provide the simplest command line required to test these files. The command must be in a single line and contain no extra text or commentary:\n"
    LintCommandPrompt = "Based on the following file list, provide the simplest command line required to lint these files. The command must be in a single line and contain no extra text or commentary:\n"
    FindReadmeCommandPrompt = `Please identify the single file from the list below that is most likely to be the README file of the project. Respond only with the file name and location, without additional commentary or explanation. List of project files:`
)
