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
  tailwindcss = pkgs.writeShellScriptBin "tailwindcss" ''
    ${pkgs.tailwindcss}/bin/tailwindcss -c ${./tailwind.config.js}
  '';
in
buildGoApplication {
  pname = "FOOdBAR";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  nativeBuildInputs = [ templ pkgs.makeWrapper tailwindcss ];
  preBuild = ''
    templ generate
    mkdir -p $out/FOOstatic
    tailwindcss -o $out/FOOstatic/tailwind.css
  '';
  postFixup = ''
    mkdir -p $out/dist
    mv $out/bin/cmd $out/dist/cmd
    makeWrapper $out/dist/cmd $out/bin/FOOdBAR \
      --set FOOdBAR_STATE ${dbpath}
  '';
  buildInputs = [ pkgs.sqlite ];
}
