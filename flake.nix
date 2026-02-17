{
  description = "A Nix-flake-based Go development environment";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1"; # unstable Nixpkgs

  outputs =
    { self, ... }@inputs:

    let
      goVersion = 24; # Change this to update the whole stack

      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forEachSupportedSystem =
        f:
        inputs.nixpkgs.lib.genAttrs supportedSystems (
          system:
          f {
            pkgs = import inputs.nixpkgs {
              inherit system;
              overlays = [ inputs.self.overlays.default ];
            };
          }
        );
    in
    {
      overlays.default = final: prev: {
        go = final."go_1_${toString goVersion}";
      };

      devShells = forEachSupportedSystem (
        { pkgs }:
        {
          default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              # go (version is specified by overlay)
              go
              # Required for .deb build
              dpkg
            ];
          };
        }
      );
      packages = forEachSupportedSystem (
        {
          pkgs,
        }:
        {
          default = pkgs.buildGoModule (finalAttrs: {
            pname = "spotiflac-cli";
            version = "2.0.0";

            src = ./.;
            vendorHash = "sha256-zU6wXQt7Vk8ks/LKx7pPmoJGBwRicUOmNI0c9byuTKI=";

            nativeBuildInputs = with pkgs; [
              installShellFiles
            ];

            subPackages = [
              "."
            ];

            postInstall = ''
              installShellCompletion --cmd spotiflac-cli \
                --bash <($out/bin/spotiflac-cli completion bash) \
                --fish <($out/bin/spotiflac-cli completion fish) \
                --zsh <($out/bin/spotiflac-cli completion zsh) 
            '';
          });
        }
      );
    };
}
