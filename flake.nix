{
  description = "A basic gomod2nix flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.inputs.flake-utils.follows = "flake-utils";
    templ.url = "github:a-h/templ";
    htmx = {
      url = "github:bigskysoftware/htmx";
      flake = false;
    };
    hyperscript = {
      url = "github:bigskysoftware/_hyperscript";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix, ... }@inputs: let
    forEachSystem = flake-utils.lib.eachSystem inputs.flake-utils.lib.allSystems;
  in
  forEachSystem (system: let
    pkgs = nixpkgs.legacyPackages.${system};

    # The current default sdk for macOS fails to compile go projects, so we use a newer one for now.
    # This has no effect on other platforms.
    callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;

    # NOTE: the program
    default = callPackage ./. {
      inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
      inherit inputs;
    };
    # NOTE: the build environment
    devShellDefault = callPackage ./shell.nix {
      inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
      inherit inputs;
    };

    /*NOTE:
      commands to build, load, and run with published port and persistent volume:
      nix build .#docker.x86_64-linux.default
      # (on zsh remember to escape #) then:
      docker load < ./result
      docker run -p 8080:8080 --mount source=foodvol,target=/var/db/foodb --rm birdee.io/foodbar
    */
    docked = pkgs.dockerTools.buildLayeredImage {
      name = "birdee.io/FOOdBAR";
      tag = "latest";
      # contents = with pkgs; [
      #   cacert
      # ];
      config = {
        Memory = 2048;
        CPUShares = 2;
        Cmd = [
          "${default}/bin/FOOdBAR"
          "-dbpath"
          "/var/db/foodb"
          "-ip"
          "0.0.0.0"
          "-port"
          "8080"
        ];
        ExposedPorts = { "8080/tcp" = {}; };
        Volumes = { "/var/db/foodb" = {}; };
      };
    };

  in
  {
    docker.default = docked;    
    packages.default = default;
    devShells.default = devShellDefault;
  }) ;
}
