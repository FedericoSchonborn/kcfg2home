{lib, buildGoModule, ...}:
buildGoModule {
  pname = "kcfg2home";
  version = "unstable";

  src = ./.;
  vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
}
