{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, buildGoApplication ? pkgs.buildGoApplication
, inputs ? {}
, dbpath ? "/tmp"
}: let
  templ = inputs.templ.packages.${pkgs.system}.templ;
  # tailwindcss = pkgs.writeShellScriptBin "tailwindcss" ''
  #   ${pkgs.tailwindcss}/bin/tailwindcss -c ${./tailwind.config.js}
  # '';
in
buildGoApplication {
  pname = "FOOdBAR";
  version = "0.1";
  pwd = ./cmd/FOOdBAR;
  src = ./.;
  modules = ./gomod2nix.toml;
  buildInputs = [ pkgs.sqlite ];
  nativeBuildInputs = [ templ pkgs.makeWrapper pkgs.tailwindcss ];
  # TODO: why does it not bundle the tailwind css?
  # The pwd and ls -l * say it is in the right place at the right time...
  preBuild = ''
    mkdir -p ./static
    tailwindcss -o ./static/tailwind.css --minify
    pwd
    ls -l *
    templ generate
    pwd
    ls -l *
  '';
  postFixup = ''
    wrapProgram $out/bin/FOOdBAR \
      --set FOOdBAR_STATE ${dbpath}
  '';
}
