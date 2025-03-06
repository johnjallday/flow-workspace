The following code structure with functions is provided:

├── LICENSE
├── README.md
├── build
│   └── build.sh
├── cmd
│   └── main.go
│       ├── [32;1mmain[0;22m()
├── docs
│   ├── flow workspace manager.png
│   ├── installation.md
│   ├── logo.webp
│   ├── requirements.md
│   └── roadmap.md
├── go.mod
├── go.sum
├── internal
│   ├── agent
│   │   └── agent.go
│   │       ├── [32;1mgetCommands[0;22m([35magentPath string[0m)[34m -> ([]Command, error)[0m
│   │       ├── [32;1mprintCommands[0;22m([35mcommands []Command[0m)
│   │       ├── [92;1mLaunchAgent[0;22m([35magentPath string[0m, [35mcommandName string[0m)
│   │       ├── [92;1mStartAgentREPL[0;22m([35magentPath string[0m)
│   ├── commands
│   │   └── commands.go
│   │       ├── [92;1mName[0;22m()[34m -> string[0m
│   │       ├── [92;1mDescription[0;22m()[34m -> string[0m
│   │       ├── [92;1mExecute[0;22m([35margs []string[0m)[34m -> error[0m
│   │       ├── [92;1mNewCommand[0;22m([35mname[0m, [35mdescription string[0m, [35mhandler func(args []string[0m)[34m -> error) Command[0m
│   ├── db
│   │   └── todo
│   │       └── todo.go
│   │           ├── [92;1mInitDB[0;22m([35mdbPath string[0m)[34m -> (*sql.DB, error)[0m
│   │           ├── [92;1mCreateTodoTable[0;22m([35mdb *sql.DB[0m)[34m -> error[0m
│   ├── project
│   │   ├── launch.go
│   │       ├── [92;1mLaunchProject[0;22m([35mprojectDir string[0m, [35msessionName string[0m)[34m -> error[0m
│   │   ├── music.go
│   │   ├── project.go
│   │       ├── [92;1mLoadProjectInfo[0;22m([35mfilename string[0m)[34m -> (*Project, error)[0m
│   │   └── repl.go
│   │       ├── [92;1mStartProjectREPL[0;22m([35mdbPath string[0m, [35mprojectDir string[0m)
│   │       ├── [32;1mprintProjectHelp[0;22m()
│   │       ├── [32;1mclearScreen[0;22m()
│   ├── repl
│   │   └── repl.go
│   │       ├── [92;1mStartREPL[0;22m([35mdbPath string[0m)
│   │       ├── [32;1mclearScreen[0;22m()
│   ├── root
│   │   ├── repl.go
│   │       ├── [92;1mStartRootREPL[0;22m([35mdbPath string[0m, [35mrootDir string[0m)
│   │       ├── [32;1mselectWorkspace[0;22m([35mrootDir string[0m, [35mreader *bufio.Reader[0m)[34m -> string[0m
│   │       ├── [32;1mprintRootHelp[0;22m()
│   │   └── root.go
│   │       ├── [92;1mListWorkspaces[0;22m([35mrootDir string[0m)
│   │       ├── [92;1mListProjects[0;22m([35mrootDir string[0m)
│   │       ├── [92;1mListAllTodos[0;22m([35mrootDir string[0m)
│   ├── startup
│   │   └── startdb.go
│   │       ├── [32;1mcreateConfig[0;22m([35mdb *sql.DB[0m, [35musername string[0m)[34m -> error[0m
│   │       ├── [32;1mcreateDB[0;22m([35mdbPath string[0m, [35musername string[0m)
│   │       ├── [92;1mStartDB[0;22m()[34m -> string[0m
│   ├── todo
│   │   ├── common.go
│   │       ├── [92;1mLoadAllTodos[0;22m([35mfilename string[0m)[34m -> ([]Todo, error)[0m
│   │       ├── [32;1mparseTodo[0;22m([35mline string[0m)[34m -> (Todo, error)[0m
│   │       ├── [92;1mWriteFileContent[0;22m([35mfilename[0m, [35mcontent string[0m)[34m -> error[0m
│   │       ├── [92;1mReadFileContent[0;22m([35mfilename string[0m)[34m -> (string, error)[0m
│   │       ├── [92;1mSaveTodos[0;22m([35mfilename string[0m, [35mtodos []Todo[0m)[34m -> error[0m
│   │   ├── migrate.go
│   │       ├── [92;1mInsertTodo[0;22m([35mdb *sql.DB[0m, [35mt Todo[0m)[34m -> error[0m
│   │       ├── [92;1mMigrateFinishedTodos[0;22m([35mtodoPath string[0m, [35mdb *sql.DB[0m)[34m -> error[0m
│   │   ├── print.go
│   │       ├── [92;1mPrintTodos[0;22m([35mtodos []Todo[0m)
│   │   ├── repl.go
│   │       ├── [92;1mStartTodoREPL[0;22m([35mdbPath string[0m, [35mtodoFilePath string[0m)
│   │       ├── [32;1mprintHelp[0;22m()
│   │       ├── [32;1mclearScreen[0;22m()
│   │   ├── service.go
│   │       ├── [92;1mNewFileTodoService[0;22m([35mtodoFilePath string[0m)[34m -> *FileTodoService[0m
│   │       ├── [92;1mAddTodo[0;22m([35mdescription[0m, [35mdueDate string[0m)[34m -> error[0m
│   │       ├── [92;1mListTodos[0;22m()[34m -> ([]Todo, error)[0m
│   │       ├── [92;1mEditTodo[0;22m([35mindex int[0m, [35mnewDescription[0m, [35mnewDueDate string[0m)[34m -> error[0m
│   │       ├── [92;1mDeleteTodo[0;22m([35mindex int[0m)[34m -> error[0m
│   │       ├── [92;1mCompleteTodo[0;22m([35mindex int[0m)[34m -> error[0m
│   │   ├── tag.go
│   │       ├── [32;1mtagProject[0;22m([35mprojectPath string[0m)[34m -> string[0m
│   │       ├── [32;1mtagWorkspace[0;22m([35mprojectPath string[0m)[34m -> string[0m
│   │       ├── [32;1mbuildTaskLine[0;22m([35mdescription[0m, [35mdueDate[0m, [35mprojectName[0m, [35mworkspaceName string[0m)[34m -> string[0m
│   │   └── todo.go
│   └── workspace
│       ├── repl.go
│           ├── [92;1mStartWorkspaceREPL[0;22m([35mdbPath string[0m, [35mworkspaceDir string[0m)
│           ├── [32;1mprintWorkspaceHelp[0;22m()
│           ├── [32;1mselectProject[0;22m([35mdbPath string[0m, [35mworkspaceDir string[0m, [35mprojs *Projects[0m, [35mreader *bufio.Reader[0m)
│       ├── tag.go
│           ├── [92;1mLoadTOML[0;22m([35mfilename string[0m)[34m -> (*WorkspaceInfo, error)[0m
│           ├── [92;1mSaveTOML[0;22m([35mfilename string[0m, [35mworkspace *WorkspaceInfo[0m)[34m -> error[0m
│       └── workspace.go
│           ├── [92;1mLoadProjectsToml[0;22m([35mworkspacePath string[0m)[34m -> (*Projects, error)[0m
│           ├── [92;1mSaveProjectsToml[0;22m([35mprojs *Projects[0m, [35mworkspacePath string[0m)[34m -> error[0m
│           ├── [92;1mUpdateProjects[0;22m([35mworkspaceDir string[0m)[34m -> (*Projects, error)[0m
│           ├── [92;1mScanAndAggregateProjects[0;22m([35mrootDir string[0m)[34m -> (*Projects, error)[0m
│           ├── [92;1mListProjects[0;22m([35mprojs *Projects[0m)
│           ├── [92;1mListAllTodos[0;22m([35mworkspaceDir string[0m)
├── project_info.toml
├── settings.toml
└── todo.md


Please implement any missing functions or suggest improvements as needed.