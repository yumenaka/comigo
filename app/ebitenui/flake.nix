{
  #TODO:Nix development dependencies for ebiten
  description = "Pure and reproducible nix overlay of binary distributed golang toolchains";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      with pkgs;
      {
        devShells.default = mkShell {
          buildInputs = [
          go
          alsa-lib
          libGL
          libX11 libXcursor libXext libXi libXinerama libXrandr
          libXxf86vm
          ];
            nativeBuildInputs = [
              go-licenses
              pkg-config
              zip
              advancecomp
              makeWrapper
            ];
        };
      }
    );
}