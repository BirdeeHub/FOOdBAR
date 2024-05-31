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
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  buildInputs = [ pkgs.sqlite ];
  nativeBuildInputs = [ templ pkgs.makeWrapper pkgs.tailwindcss ];
  postUnpack = ''
    targetStaticDir=$TEMPDIR/$sourceRoot/static
    mkdir -p $targetStaticDir
    tailwindcss -o $targetStaticDir/tailwind.css -c ${./tailwind.config.js} --minify
    cp ${inputs.htmx}/dist/htmx.min.js $targetStaticDir
    cp ${inputs.hyperscript}/dist/_hyperscript.min.js $targetStaticDir
  '';
  preBuild = ''
    templ generate
  '';
  postFixup = ''
    # https://github.com/NixOS/nixpkgs/blob/master/pkgs/build-support/setup-hooks/make-wrapper.sh
    wrapProgram $out/bin/FOOdBAR \
      --set FOOdBAR_STATE ${dbpath}
  '';
}
