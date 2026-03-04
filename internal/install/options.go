package install

type Options struct {
	Minimal        bool
	NoPackages     bool
	SkipOptional   bool
	Profile        string
	DotfilesDir    string
	NonInteractive bool
	DryRun         bool
	Verbose        bool
}
