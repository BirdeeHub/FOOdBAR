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
  # It is in the right place before build occurs...
  preBuild = ''
    mkdir -p ./static
    templ generate
    tailwindcss -o ./static/tailwind.css --minify -c ${./tailwind.config.js}
    pwd
    ls -l *
    cat ./static/tailwind.css
  '';
  postFixup = ''
    wrapProgram $out/bin/FOOdBAR \
      --set FOOdBAR_STATE ${dbpath}
  '';
}
