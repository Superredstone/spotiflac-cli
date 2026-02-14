{
  description = "A Nix-flake-based Go development environment";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1"; # unstable Nixpkgs

  outputs =
    { self, ... }@inputs:

    let
      pname = "spotiflac-cli";
      version = "7.0.9";
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

      pkgs = import inputs.nixpkgs {
        system = "x86_64-linux";
        overlays = [ inputs.self.overlays.default ];
      };

      spotiflac = pkgs.fetchFromGitHub {
        owner = "afkarxyz";
        repo = "SpotiFLAC";
        tag = "v${version}";
        hash = "sha256-VHYof17C+eRoZfssXRQpbB8GXlcfPhyRiWltM6yDqe0=";
      };
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

              ffmpeg_7
            ];
          };
        }
      );
      packages = forEachSupportedSystem (
        { pkgs }:
        {
          default = pkgs.buildGoModule (finalAttrs: {
            inherit pname version;
            src = ./.;
            vendorHash = "sha256-EpGgfiCqJjHEOphV2x8FmXeIFls7eq2NVxb/or4NLUo=";

            dependencies = with pkgs; [
              ffmpeg_7
            ];

            nativeBuildInputs = with pkgs; [
              installShellFiles
            ];

            subPackages = [
              "."
            ];

            postPatch = ''
              cp -r ${spotiflac} ./SpotiFLAC/
              sed -i "s/git clone https:\/\/github.com\/afkarxyz\/SpotiFLAC.git//g" ./tools/fetch_spotiflac_backend.sh
              sed -i "s/rm -rf SpotiFLAC//g" ./tools/fetch_spotiflac_backend.sh
              ./tools/fetch_spotiflac_backend.sh
            '';

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
