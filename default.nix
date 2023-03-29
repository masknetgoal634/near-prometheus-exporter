with import <nixpkgs> {};
mkShell {
  nativeBuildInputs = [
    go
    bashInteractive
  ];
}
