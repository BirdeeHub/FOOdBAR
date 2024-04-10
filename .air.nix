{ pkgs, templ, go, gomod2nix, ... }:
pkgs.writeText "aircfg" (/*toml*/''
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  # bin = "./result/bin/cmd"
  # pre_cmd = [ "git add ." ]
  # cmd ="nix build --show-trace"
  bin = "./tmp/main"
  pre_cmd = [ "${go}/bin/go mod tidy", "${gomod2nix}/bin/gomod2nix" ]
  cmd = "${templ}/bin/templ generate && ${go}/bin/go build -o ./tmp/main cmd/*.go"
  delay = 0
  exclude_dir = ["node_modules", "assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go", ".*_templ.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = ["cmd", "views", "pkg"]
  include_ext = ["go", "html", "tpl", "tmpl", "templ"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
'')
