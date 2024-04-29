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
  nativeBuildInputs = with pkgs; [ templ makeWrapper ];
  preBuild = ''
    templ generate
  '';
  postFixup = ''
    mkdir -p $out/dist
    mv $out/bin/cmd $out/dist/cmd
    makeWrapper $out/dist/cmd $out/bin/FOOdBAR \
      --set FOOdBAR_STATE ${dbpath}
  '';
  buildInputs = [ pkgs.sqlite ];
}
