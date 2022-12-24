{
  description = "Generate NixOS/Home Manager modules from KDE Frameworks configuration files";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = {nixpkgs, ...}: let
    forAllSystems = f: nixpkgs.lib.genAttrs ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"] (system: f nixpkgs.legacyPackages.${system});
  in {
    packages = forAllSystems (pkgs: rec {
      kcfg2home = pkgs.callPackage ./. {};
      default = kcfg2home;
    });

    devShells = forAllSystems (pkgs: {
      default = import ./shell.nix {inherit pkgs;};
    });

    formatter = forAllSystems (pkgs: pkgs.alejandra);
  };
}
