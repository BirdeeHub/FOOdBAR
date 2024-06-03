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
, dbpath ? "''"
, keypath ? "''"
, ip ? "''"
, port ? "''"
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
    newsource=$TEMPDIR/''${sourceRoot}
    targetStaticDir=$newsource/static
    targetFOOstaticDir=$newsource/FOOstatic
    mkdir -p $targetStaticDir
    mkdir -p $targetFOOstaticDir

    echo "generate tailwind.css"
    tailwindcss -o $targetStaticDir/tailwind.css -c ${./tailwind.config.js} --minify

    echo "gzipping select files"
    gzip -k -c $targetFOOstaticDir/foodbarfavicon.svg > $targetFOOstaticDir/foodbarfavicon.svg.gz
    gzip -k -c $targetStaticDir/foodbarloginfavicon.svg > $targetStaticDir/foodbarloginfavicon.svg.gz
    gzip -k -c $targetStaticDir/tailwind.css > $targetStaticDir/tailwind.css.gz
    gzip -k -c $targetStaticDir/patchhyperscript.js > $targetStaticDir/patchhyperscript.js.gz

    echo "bundling client side dependencies"
    gzip -k -c ${inputs.htmx}/dist/htmx.min.js > $targetStaticDir/htmx.min.js.gz
    gzip -k -c ${inputs.hyperscript}/dist/_hyperscript.min.js > $targetStaticDir/_hyperscript.min.js.gz
    ls -l $targetStaticDir/* $targetFOOstaticDir/*
  '';
  preBuild = ''
    templ generate
  '';
  postFixup = ''
    # https://github.com/NixOS/nixpkgs/blob/master/pkgs/build-support/setup-hooks/make-wrapper.sh
    # wrapProgram $out/bin/FOOdBAR \
    #   --set FOOdBAR_STATE ${dbpath}\
    #   --set FOOdBAR_SIGNING_KEY ${keypath}\
    #   --set FOOdBAR_ADDRESS ${ip}\
    #   --set FOOdBAR_PORT ${port}\
    #   --add-flags "-port ${port} -ip ${ip} -dbpath ${dbpath} -keypath ${keypath}"
  '';
}
