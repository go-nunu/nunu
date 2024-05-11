package config

var (
	Version       = "1.1.0"
	WireCmd       = "github.com/google/wire/cmd/wire@latest"
	NunuCmd       = "github.com/go-nunu/nunu@latest"
	RepoBase      = "https://github.com/go-nunu/nunu-layout-base.git"
	RepoAdvanced  = "https://github.com/go-nunu/nunu-layout-advanced.git"
	RepoChat      = "https://github.com/go-nunu/nunu-layout-chat.git"
	RunExcludeDir = ".git,.idea,tmp,vendor"
	RunIncludeExt = "go,html,yaml,yml,toml,ini,json,xml,tpl,tmpl"
)
