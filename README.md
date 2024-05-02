
## ephemyral

Ephemyral is an AI-powered CLI application for managing coding tasks that leverage machine learning models.

### Synopsis

                    █████                                                              ████ 
                   ░░███                                                              ░░███ 
  ██████  ████████  ░███████    ██████  █████████████   █████ ████ ████████   ██████   ░███ 
 ███░░███░░███░░███ ░███░░███  ███░░███░░███░░███░░███ ░░███ ░███ ░░███░░███ ░░░░░███  ░███ 
░███████  ░███ ░███ ░███ ░███ ░███████  ░███ ░███ ░███  ░███ ░███  ░███ ░░░   ███████  ░███ 
░███░░░   ░███ ░███ ░███ ░███ ░███░░░   ░███ ░███ ░███  ░███ ░███  ░███      ███░░███  ░███ 
░░██████  ░███████  ████ █████░░██████  █████░███ █████ ░░███████  █████    ░░████████ █████
 ░░░░░░   ░███░░░  ░░░░ ░░░░░  ░░░░░░  ░░░░░ ░░░ ░░░░░   ░░░░░███ ░░░░░      ░░░░░░░░ ░░░░░ 
          ░███                                           ███ ░███                           
          █████                                         ░░██████                            
         ░░░░░                                           ░░░░░░                             

Ephemyral is an AI-powered CLI application designed to streamline and optimize various software development tasks with the help of machine learning. By leveraging large language models, Ephemyral provides a set of robust commands that simplify building, testing, and managing development workflows. This tool is tailored for software engineers, data scientists, and anyone managing software projects.

### Options

```
      --config string   config file (default is $HOME/.ephemyral.yaml)
  -h, --help            help for ephemyral
```

### SEE ALSO

* [ephemyral build](./docs/ephemyral_build.md)	 - Use AI to intelligently generate and execute a build commands for the specified directory, optimizing for performance and efficiency.
* [ephemyral create](./docs/ephemyral_create.md)	 - Employ a language model to generate new code files based on a natural language prompt. If the file path is a directory, it generates multiple AI-crafted files.
* [ephemyral docs](./docs/ephemyral_docs.md)	 - Generate and execute commands to create documentation, enhancing your codebase's maintainability.
* [ephemyral init](./docs/ephemyral_init.md)	 - Initialize a new .ephemyral file for AI-generated configurations, setting up the environment in a directory.
* [ephemyral lint](./docs/ephemyral_lint.md)	 - Use machine learning models to generate and execute a lint commands, improving code quality by identifying patterns and anomalies.
* [ephemyral refactor](./docs/ephemyral_refactor.md)	 - Utilize an advanced LLM to refactor a give file or all files in a provided directory based on prompts, outputting, building and testing the improved code.
* [ephemyral test](./docs/ephemyral_test.md)	 - Deploy AI models to generate and run optimized test commands for the specified directories, enhancing test accuracy and efficiency.


## ephemyral build

Use AI to intelligently generate and execute a build commands for the specified directory, optimizing for performance and efficiency.

### Synopsis

The 'build' command generates a building command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new build command and executes it. Use this command to ensure your project builds correctly and is free from errors.

```
ephemyral build [directory] [flags]
```

### Options

```
  -h, --help        help for build
      --retry int   Number of retries for generating and executing build command (default 3)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.ephemyral.yaml)
```

### SEE ALSO

* [ephemyral](./docs/ephemyral.md)	 - Ephemyral is an AI-powered CLI application for managing coding tasks that leverage machine learning models.


## ephemyral create

Employ a language model to generate new code files based on a natural language prompt. If the file path is a directory, it generates multiple AI-crafted files.

### Synopsis

This command generates a new code file based on a given prompt. 
If the file path is a directory, it uses a query to determine the file names and creates new files based on the provided prompt.

```
ephemyral create [file path] [prompt] [flags]
```

### Options

```
  -h, --help   help for create
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.ephemyral.yaml)
```

### SEE ALSO

* [ephemyral](./docs/ephemyral.md)	 - Ephemyral is an AI-powered CLI application for managing coding tasks that leverage machine learning models.


## ephemyral docs

Generate and execute commands to create documentation, enhancing your codebase's maintainability.

### Synopsis

The 'docs' command creates a command to produce documentation (like a README or API documentation) for the project's files. It then updates the '.ephemyral' configuration file with the new command and executes it.

```
ephemyral docs [directory] [flags]
```

### Options

```
  -h, --help        help for docs
      --retry int   Number of retries for generating and executing docs command (default 3)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.ephemyral.yaml)
```

### SEE ALSO

* [ephemyral](./docs/ephemyral.md)	 - Ephemyral is an AI-powered CLI application for managing coding tasks that leverage machine learning models.


## ephemyral init

Initialize a new .ephemyral file for AI-generated configurations, setting up the environment in a directory.

### Synopsis

The 'init' command is designed to help set up the necessary configurations for an Ephemyral-based project. When executed, the command checks if there is an existing '.ephemyral' file in the current directory. This file contains YAML-formatted information related to AI-generated build, test, and lint commands.
If the '.ephemyral' file does not exist, the command creates one with a basic template for build, test, and lint command configurations. This is useful for initializing a new Ephemyral task or project where AI-driven commands can be defined and customized later.
If a '.ephemyral' file is already present, the command confirms that the Ephemyral task has been initialized, allowing users to proceed with other tasks such as building, testing, or linting.
The newly created '.ephemyral' file has a default structure with placeholders for build, test, and lint commands, which can be edited as needed. The 'init' command provides a foundation for AI-based project management, ensuring that an essential configuration file is in place before additional tasks are performed.

```
ephemyral init [flags]
```

### Options

```
  -h, --help   help for init
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.ephemyral.yaml)
```

### SEE ALSO

* [ephemyral](./docs/ephemyral.md)	 - Ephemyral is an AI-powered CLI application for managing coding tasks that leverage machine learning models.


## ephemyral lint

Use machine learning models to generate and execute a lint commands, improving code quality by identifying patterns and anomalies.

### Synopsis

The 'lint' command generates a linting command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new linting command and executes it. Use this command to ensure your project adheres to coding standards and is free from basic syntax errors.

```
ephemyral lint [directory] [flags]
```

### Options

```
  -h, --help        help for lint
      --retry int   Number of retries for generating and executing lint command (default 3)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.ephemyral.yaml)
```

### SEE ALSO

* [ephemyral](./docs/ephemyral.md)	 - Ephemyral is an AI-powered CLI application for managing coding tasks that leverage machine learning models.


## ephemyral refactor

Utilize an advanced LLM to refactor a give file or all files in a provided directory based on prompts, outputting, building and testing the improved code.

### Synopsis

This command refactors a given file or all files in a directory by sending a prompt to an LLM 
and applying the suggested changes, replacing the file content or creating new files in the specified directory.

```
ephemyral refactor [file path] [prompt] [new file path] [flags]
```

### Options

```
  -h, --help        help for refactor
      --retry int   Number of retries for refactoring files (default 3)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.ephemyral.yaml)
```

### SEE ALSO

* [ephemyral](./docs/ephemyral.md)	 - Ephemyral is an AI-powered CLI application for managing coding tasks that leverage machine learning models.


## ephemyral test

Deploy AI models to generate and run optimized test commands for the specified directories, enhancing test accuracy and efficiency.

### Synopsis

The 'test' command generates a testing command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new testing command and executes it. Use this command to ensure your project adheres to testing standards and is free from test errors.

```
ephemyral test [directory] [flags]
```

### Options

```
  -h, --help        help for test
      --retry int   Number of retries for generating and executing test command (default 3)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.ephemyral.yaml)
```

### SEE ALSO

* [ephemyral](./docs/ephemyral.md)	 - Ephemyral is an AI-powered CLI application for managing coding tasks that leverage machine learning models.


