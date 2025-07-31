{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  packages = with pkgs; [
    act
    go
    golangci-lint
    gotestsum
  ];

  shellHook = ''
    export GOPATH=$(pwd)/.go
  '';
}
