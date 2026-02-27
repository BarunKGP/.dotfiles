package install

type Options struct {
	Minimal        bool
	NoPackages     bool
	SkipOptional   bool
	DotfilesDir    string
	NonInteractive bool
	DryRun         bool
	Verbose        bool
}
