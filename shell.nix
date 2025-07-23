{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  packages = with pkgs; [
    go
    golangci-lint
    gotestsum
  ];

  shellHook = ''
    export GOPATH=$(pwd)/.go
  '';
}
