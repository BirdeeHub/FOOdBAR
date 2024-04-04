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
, mkGoEnv ? pkgs.mkGoEnv
, gomod2nix ? pkgs.gomod2nix
, inputs ? {}
}:

let
  templ = inputs.templ.packages.${pkgs.system}.templ;
  air = pkgs.writeShellScriptBin "air" ''
    ${pkgs.air}/bin/air -c ${import ./.air.nix { inherit pkgs air templ gomod2nix goEnv; inherit (pkgs) go; }}
  '';
  goEnv = mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  DEVSHELL = 0;
  packages = [
    goEnv
    gomod2nix
    air
    templ
  ];
  shellHook = ''
    exec ${pkgs.zsh}/bin/zsh
  '';
}
