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
  pwd = ./cmd/FOOdBAR;
  src = ./.;
  modules = ./gomod2nix.toml;
  nativeBuildInputs = [ templ pkgs.makeWrapper tailwindcss ];
  preBuild = ''
    templ generate
    tailwindcss build > ./static/tailwind.css
  '';
  postFixup = ''
    wrapProgram $out/bin/FOOdBAR \
      --set FOOdBAR_STATE ${dbpath}
  '';
  buildInputs = [ pkgs.sqlite ];
}
